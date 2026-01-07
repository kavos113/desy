package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity

@Entity(
  tableName = "related_courses",
  primaryKeys = ["lecture_id", "related_lecture_id"],
)
data class RelatedCourseEntity(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "related_lecture_id")
  val relatedLectureId: Int,
)
