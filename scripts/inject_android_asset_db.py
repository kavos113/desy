#!/usr/bin/env python3
"""Inject a SQLite DB file into Android assets.

Copies the specified database file into:
  android/app/src/main/assets/dasy_database.db

Usage:
  python scripts/inject_android_asset_db.py path/to/source.db
"""

from __future__ import annotations

import argparse
import os
import shutil
import sys

SQLITE_HEADER = b"SQLite format 3\x00"


def _is_sqlite_file(path: str) -> bool:
  try:
    with open(path, "rb") as f:
      header = f.read(len(SQLITE_HEADER))
    return header == SQLITE_HEADER
  except OSError:
    return False


def main(argv: list[str]) -> int:
  parser = argparse.ArgumentParser(
    description=(
      "コマンドライン引数で指定したSQLite DBを、"
      "android/app/src/main/assets/dasy_database.db に注入(置き換えコピー)します。"
    )
  )
  parser.add_argument(
    "source_db",
    help="注入元のSQLite DBファイルパス",
  )
  args = parser.parse_args(argv)

  repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))
  src = os.path.abspath(args.source_db)
  dst = os.path.join(repo_root, "android", "app", "src", "main", "assets", "dasy_database.db")

  if not os.path.isfile(src):
    print(f"ERROR: 注入元DBが見つかりません: {src}", file=sys.stderr)
    return 2

  if not _is_sqlite_file(src):
    print(f"ERROR: SQLiteファイルではありません: {src}", file=sys.stderr)
    return 3

  os.makedirs(os.path.dirname(dst), exist_ok=True)

  try:
    shutil.copy2(src, dst)
  except OSError as e:
    print(f"ERROR: 注入に失敗しました: {e}", file=sys.stderr)
    return 4

  size = os.path.getsize(dst)
  print(f"OK: {dst} に注入しました ({size} bytes)")
  return 0


if __name__ == "__main__":
  raise SystemExit(main(sys.argv[1:]))
