package com.github.kavos113.desy.ui.detail

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.Teacher
import com.github.kavos113.desy.domain.TimeTable

internal fun splitIntoLines(value: String?): List<String> {
  return value
    ?.split(Regex("\\r?\\n"))
    ?.map { it.trimEnd() }
    ?.filter { it.isNotBlank() }
    ?: emptyList()
}

internal fun formatTeachers(teachers: List<Teacher>): String {
  return teachers
    .mapNotNull { it.name?.trim()?.takeIf { name -> name.isNotEmpty() } }
    .distinct()
    .joinToString(", ")
}

internal fun formatTimetablesWithRoom(timetables: List<TimeTable>): String {
  if (timetables.isEmpty()) return ""

  return timetables
    .sortedWith(
      compareBy<TimeTable>(
        { it.semester?.ordinal ?: Int.MAX_VALUE },
        { it.dayOfWeek?.ordinal ?: Int.MAX_VALUE },
        { it.period ?: Int.MAX_VALUE },
      )
    )
    .joinToString(", ") { tt ->
      val day = tt.dayOfWeek?.toJapaneseLabel().orEmpty()
      val period = tt.period?.takeIf { it > 0 }?.toString().orEmpty()
      val room = tt.room?.name?.trim().orEmpty()
      val core = (day + period).ifBlank { "?" }
      if (room.isBlank()) core else "$core($room)"
    }
}

internal fun formatSemesters(timetables: List<TimeTable>): String {
  val semesters = timetables.mapNotNull { it.semester }.distinct()
  if (semesters.isEmpty()) return ""
  return semesters.joinToString("/") { it.toJapaneseLabel() }
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
