package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "lecture_plans")
data class LecturePlanEntity(
  @PrimaryKey(autoGenerate = true)
  @ColumnInfo(name = "id")
  val id: Int = 0,

  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "count")
  val count: Int?,

  @ColumnInfo(name = "plan")
  val plan: String?,

  @ColumnInfo(name = "assignment")
  val assignment: String?,
)
