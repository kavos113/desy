package com.github.kavos113.desy.repository.db.model

import androidx.room.ColumnInfo

/** Roomのクエリ結果を受けるための軽量DTO */
data class TeacherRow(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,
  @ColumnInfo(name = "id")
  val teacherId: Int,
  @ColumnInfo(name = "name")
  val name: String,
  @ColumnInfo(name = "url")
  val url: String?,
)
