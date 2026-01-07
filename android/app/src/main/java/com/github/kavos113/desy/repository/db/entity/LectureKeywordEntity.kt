package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity

@Entity(
  tableName = "lecture_keywords",
  primaryKeys = ["lecture_id", "keyword"],
)
data class LectureKeywordEntity(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "keyword")
  val keyword: String,
)
