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
      compareBy<TimeTable>(
        { it.semester?.ordinal ?: Int.MAX_VALUE },
        { it.dayOfWeek?.ordinal ?: Int.MAX_VALUE },
        { it.period ?: Int.MAX_VALUE },
      )
    )
    .joinToString(",") { tt ->
      val day = tt.dayOfWeek?.toJapaneseLabel().orEmpty()
      val period = tt.period?.takeIf { it > 0 }?.toString().orEmpty()
      val core = (day + period).ifBlank { "?" }
      val semester = tt.semester?.toJapaneseLabel()
      if (semester.isNullOrBlank()) core else "$core($semester)"
    }
}

private fun Semester.toJapaneseLabel(): String = when (this) {
  Semester.spring -> "春"
  Semester.summer -> "夏"
  Semester.fall -> "秋"
  Semester.winter -> "冬"
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
