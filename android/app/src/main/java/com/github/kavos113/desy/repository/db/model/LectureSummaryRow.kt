package com.github.kavos113.desy.repository.db.model

import androidx.room.ColumnInfo

/** Roomのクエリ結果を受けるための軽量DTO */
data class LectureSummaryRow(
  @ColumnInfo(name = "id")
  val id: Int,
  @ColumnInfo(name = "university")
  val university: String,
  @ColumnInfo(name = "title")
  val title: String,
  @ColumnInfo(name = "department")
  val department: String,
  @ColumnInfo(name = "code")
  val code: String,
  @ColumnInfo(name = "level")
  val level: Int?,
  @ColumnInfo(name = "credit")
  val credit: Int?,
  @ColumnInfo(name = "year")
  val year: Int?,
)
