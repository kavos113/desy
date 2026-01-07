package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity

@Entity(
  tableName = "lecture_teachers",
  primaryKeys = ["lecture_id", "teacher_id"],
)
data class LectureTeacherCrossRef(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "teacher_id")
  val teacherId: Int,
)
