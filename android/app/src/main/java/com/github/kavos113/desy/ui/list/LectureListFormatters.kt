package com.github.kavos113.desy.ui.list

import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.ui.formatters.TimeTableFormatter

internal fun formatLectureOpenTerm(timetables: List<TimeTable>): String {
  val semesters = timetables.mapNotNull { it.semester }.distinct()
  if (semesters.isEmpty()) return ""
  return semesters.joinToString("/") { semester -> semester.toJapaneseLabel() }
}

internal fun formatLectureTimetable(timetables: List<TimeTable>): String {
  return TimeTableFormatter.format(timetables = timetables, includeRoom = false)
}

private fun Semester.toJapaneseLabel(): String = when (this) {
  Semester.spring -> "1Q"
  Semester.summer -> "2Q"
  Semester.fall -> "3Q"
  Semester.winter -> "4Q"
}

// 時間割の表示は TimeTableFormatter に集約