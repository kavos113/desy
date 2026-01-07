package com.github.kavos113.desy.repository.db.model

import androidx.room.ColumnInfo

/** Roomのクエリ結果を受けるための軽量DTO */
data class TimetableRow(
  @ColumnInfo(name = "lecture_id")
  val lectureId: Int,
  @ColumnInfo(name = "semester")
  val semester: String?,
  @ColumnInfo(name = "room_id")
  val roomId: Int?,
  @ColumnInfo(name = "room_name")
  val roomName: String?,
  @ColumnInfo(name = "day_of_week")
  val dayOfWeek: String?,
  @ColumnInfo(name = "period")
  val period: Int?,
)
