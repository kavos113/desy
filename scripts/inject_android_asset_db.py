#!/usr/bin/env python3
"""Inject a SQLite DB file into Android assets.

Copies data (SELECT & INSERT) from the specified database file into:
  android/app/src/main/assets/dasy_database.db

Usage:
  python scripts/inject_android_asset_db.py path/to/source.db
"""

from __future__ import annotations

import argparse
import os
import sqlite3
import sys
from typing import Iterable

SQLITE_HEADER = b"SQLite format 3\x00"


def _is_sqlite_file(path: str) -> bool:
  try:
    with open(path, "rb") as f:
      header = f.read(len(SQLITE_HEADER))
    return header == SQLITE_HEADER
  except OSError:
    return False


def _q(identifier: str) -> str:
  return '"' + identifier.replace('"', '""') + '"'


def _list_tables(conn: sqlite3.Connection) -> list[str]:
  rows = conn.execute(
    "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
  ).fetchall()
  return [row[0] for row in rows]


def _table_columns(conn: sqlite3.Connection, table: str) -> list[str]:
  rows = conn.execute(f"PRAGMA table_info({_q(table)})").fetchall()
  # columns: cid, name, type, notnull, dflt_value, pk
  return [row[1] for row in rows]


def _iter_rows_in_chunks(
  conn: sqlite3.Connection,
  sql: str,
  chunk_size: int,
) -> Iterable[list[tuple]]:
  cur = conn.execute(sql)
  while True:
    chunk = cur.fetchmany(chunk_size)
    if not chunk:
      break
    yield chunk


def _ensure_schema_from_source(src_conn: sqlite3.Connection, dst_conn: sqlite3.Connection) -> None:
  # Create schema objects in a deterministic order.
  rows = src_conn.execute(
    """
    SELECT type, name, sql
    FROM sqlite_master
    WHERE sql IS NOT NULL
      AND name NOT LIKE 'sqlite_%'
      AND type IN ('table', 'index', 'trigger', 'view')
    ORDER BY
      CASE type
        WHEN 'table' THEN 1
        WHEN 'view' THEN 2
        WHEN 'index' THEN 3
        WHEN 'trigger' THEN 4
        ELSE 5
      END,
      name
    """
  ).fetchall()

  for obj_type, _name, sql in rows:
    # Skip internal auto indices: sql is NULL already, but keep safe.
    if not sql:
      continue
    dst_conn.execute(sql)


def main(argv: list[str]) -> int:
  parser = argparse.ArgumentParser(
    description=(
      "コマンドライン引数で指定したSQLite DBを、"
      "android/app/src/main/assets/dasy_database.db に注入(SELECT & INSERTで中身を反映)します。"
    )
  )
  parser.add_argument(
    "source_db",
    help="注入元のSQLite DBファイルパス",
  )
  parser.add_argument(
    "--dest",
    default=None,
    help="注入先DB(デフォルト: android/app/src/main/assets/dasy_database.db)",
  )
  args = parser.parse_args(argv)

  repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))
  src = os.path.abspath(args.source_db)
  dst = (
    os.path.abspath(args.dest)
    if args.dest
    else os.path.join(repo_root, "android", "app", "src", "main", "assets", "dasy_database.db")
  )

  if not os.path.isfile(src):
    print(f"ERROR: 注入元DBが見つかりません: {src}", file=sys.stderr)
    return 2

  if not _is_sqlite_file(src):
    print(f"ERROR: SQLiteファイルではありません: {src}", file=sys.stderr)
    return 3

  os.makedirs(os.path.dirname(dst), exist_ok=True)

  src_conn = sqlite3.connect(src)
  try:
    dst_exists = os.path.isfile(dst)
    dst_conn = sqlite3.connect(dst)
  except sqlite3.Error as e:
    print(f"ERROR: 注入先DBを開けません: {dst} ({e})", file=sys.stderr)
    return 4

  try:
    # If destination DB didn't exist, create schema from source.
    if not dst_exists:
      _ensure_schema_from_source(src_conn, dst_conn)

    src_tables = set(_list_tables(src_conn))
    dst_tables = _list_tables(dst_conn)
    common_tables = [t for t in dst_tables if t in src_tables]

    if not common_tables:
      print("ERROR: 注入対象の共通テーブルが見つかりません。スキーマが一致しているか確認してください。", file=sys.stderr)
      return 5

    dst_conn.execute("PRAGMA foreign_keys=OFF")
    dst_conn.execute("BEGIN")

    for table in common_tables:
      dst_cols = _table_columns(dst_conn, table)
      src_cols = set(_table_columns(src_conn, table))

      missing = [c for c in dst_cols if c not in src_cols]
      if missing:
        raise RuntimeError(f"テーブル {table} のカラムが一致しません。注入先にのみ存在: {missing}")

      col_sql = ", ".join(_q(c) for c in dst_cols)
      placeholders = ", ".join(["?"] * len(dst_cols))

      dst_conn.execute(f"DELETE FROM {_q(table)}")

      select_sql = f"SELECT {col_sql} FROM {_q(table)}"
      insert_sql = f"INSERT INTO {_q(table)} ({col_sql}) VALUES ({placeholders})"

      for chunk in _iter_rows_in_chunks(src_conn, select_sql, chunk_size=1000):
        dst_conn.executemany(insert_sql, chunk)

    # Copy sqlite_sequence if present in both, to preserve AUTOINCREMENT counters.
    try:
      has_seq_src = src_conn.execute(
        "SELECT 1 FROM sqlite_master WHERE type='table' AND name='sqlite_sequence'"
      ).fetchone() is not None
      has_seq_dst = dst_conn.execute(
        "SELECT 1 FROM sqlite_master WHERE type='table' AND name='sqlite_sequence'"
      ).fetchone() is not None
      if has_seq_src and has_seq_dst:
        dst_conn.execute("DELETE FROM sqlite_sequence")
        rows = src_conn.execute("SELECT name, seq FROM sqlite_sequence").fetchall()
        dst_conn.executemany("INSERT INTO sqlite_sequence(name, seq) VALUES (?, ?)", rows)
    except sqlite3.Error:
      # Best-effort only.
      pass

    dst_conn.execute("COMMIT")
    dst_conn.execute("PRAGMA foreign_keys=ON")
  except Exception as e:
    try:
      dst_conn.execute("ROLLBACK")
    except sqlite3.Error:
      pass
    print(f"ERROR: 注入に失敗しました: {e}", file=sys.stderr)
    return 6
  finally:
    try:
      src_conn.close()
    except sqlite3.Error:
      pass
    try:
      dst_conn.close()
    except sqlite3.Error:
      pass

  size = os.path.getsize(dst)
  print(f"OK: {dst} に注入しました ({size} bytes)")
  return 0


if __name__ == "__main__":
  raise SystemExit(main(sys.argv[1:]))
