package com.github.kavos113.desy.domain

data class LectureSummary(
  val id: Int,
  val university: String,
  val title: String,
  val department: String? = null,
  val code: String? = null,
  val level: Level? = null,
  val credit: Int? = null,
  val year: Int? = null,
  val timetables: List<TimeTable> = emptyList(),
  val teachers: List<Teacher> = emptyList(),
)
