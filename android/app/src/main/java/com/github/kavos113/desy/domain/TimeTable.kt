package com.github.kavos113.desy.domain

data class TimeTable(
  val lectureId: Int,
  val semester: Semester? = null,
  val room: Room? = null,
  val dayOfWeek: DayOfWeek? = null,
  val period: Int? = null,
)
