package com.github.kavos113.desy.ui.detail

import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.Teacher
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.ui.formatters.TimeTableFormatter

internal fun splitIntoLines(value: String?): List<String> {
  return value
    ?.split(Regex("\\r?\\n"))
    ?.map { it.trimEnd() }
    ?.filter { it.isNotBlank() }
    ?: emptyList()
}

internal fun formatTeachers(teachers: List<Teacher>): String {
  return teachers
    .mapNotNull { it.name.trim().takeIf { name -> name.isNotEmpty() } }
    .distinct()
    .joinToString(", ")
}

internal fun formatTimetablesWithRoom(timetables: List<TimeTable>): String {
  return TimeTableFormatter.format(timetables = timetables, includeRoom = true)
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

// DayOfWeek 表示は TimeTableFormatter に集約
