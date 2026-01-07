package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity

@Entity(
  tableName = "related_course_codes",
  primaryKeys = ["lecture_id", "code"],
)
data class RelatedCourseCodeEntity(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "code")
  val code: String,
)
