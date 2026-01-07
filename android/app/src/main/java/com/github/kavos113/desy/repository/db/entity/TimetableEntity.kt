package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "timetables")
data class TimetableEntity(
  @PrimaryKey(autoGenerate = true)
  @ColumnInfo(name = "id")
  val id: Int = 0,

  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,

  @ColumnInfo(name = "semester")
  val semester: String?,

  @ColumnInfo(name = "room_id")
  val roomId: Int?,

  @ColumnInfo(name = "day_of_week")
  val dayOfWeek: String?,

  @ColumnInfo(name = "period")
  val period: Int?,
)
