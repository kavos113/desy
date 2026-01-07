package com.github.kavos113.desy.ui.list

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.TimeTable

internal fun formatLectureOpenTerm(timetables: List<TimeTable>): String {
  val semesters = timetables.mapNotNull { it.semester }.distinct()
  if (semesters.isEmpty()) return ""
  return semesters.joinToString("/") { semester -> semester.toJapaneseLabel() }
}

internal fun formatLectureTimetable(timetables: List<TimeTable>): String {
  if (timetables.isEmpty()) return ""

  return timetables
    .sortedWith(
      compareBy(
        { it.semester?.ordinal ?: Int.MAX_VALUE },
        { it.dayOfWeek?.ordinal ?: Int.MAX_VALUE },
        { it.period ?: Int.MAX_VALUE },
      )
    )
    .joinToString(",") { tt ->
      val day = tt.dayOfWeek?.toJapaneseLabel().orEmpty()
      val period = tt.period?.takeIf { it > 0 }?.toTimeTableLabel().orEmpty()
      (day + period).ifBlank { "?" }
    }
}

private fun Semester.toJapaneseLabel(): String = when (this) {
  Semester.spring -> "1Q"
  Semester.summer -> "2Q"
  Semester.fall -> "3Q"
  Semester.winter -> "4Q"
}

private fun DayOfWeek.toJapaneseLabel(): String = when (this) {
  DayOfWeek.monday -> "月"
  DayOfWeek.tuesday -> "火"
  DayOfWeek.wednesday -> "水"
  DayOfWeek.thursday -> "木"
  DayOfWeek.friday -> "金"
  DayOfWeek.saturday -> "土"
  DayOfWeek.sunday -> "日"
}

private fun Int.toTimeTableLabel(): String = when (this) {
  1 -> "1-2"
  2 -> "3-4"
  3 -> "5-6"
  4 -> "7-8"
  5 -> "9-10"
  else -> ""
}