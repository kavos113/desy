package com.github.kavos113.desy.ui.formatters

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.TimeTable

internal object TimeTableFormatter {
  private val dayOrder = listOf(
    DayOfWeek.monday,
    DayOfWeek.tuesday,
    DayOfWeek.wednesday,
    DayOfWeek.thursday,
    DayOfWeek.friday,
    DayOfWeek.saturday,
    DayOfWeek.sunday,
  )

  fun format(
    timetables: List<TimeTable>,
    includeRoom: Boolean = true,
  ): String {
    if (timetables.isEmpty()) return ""

    val groups = linkedMapOf<String, TimeTableGroup>()
    val fallbacks = mutableListOf<String>()

    timetables.forEach { timetable ->
      val day = timetable.dayOfWeek
      val dayLabel = day?.toJapaneseLabel().orEmpty()
      val roomLabel = timetable.room?.name?.trim().orEmpty()
      val periodValue = timetable.period
      val hasPeriod = periodValue != null && periodValue > 0

      val fallback = buildFallbackLabel(
        dayLabel = dayLabel,
        period = periodValue,
        roomName = timetable.room?.name,
        includeRoom = includeRoom,
      )

      if (day == null || dayLabel.isBlank() || !hasPeriod) {
        if (fallback.isNotBlank()) fallbacks.add(fallback)
        return@forEach
      }

      val key = "${day.name}::${roomLabel}"
      val group = groups.getOrPut(key) {
        TimeTableGroup(
          day = day,
          dayLabel = dayLabel,
          roomLabel = roomLabel,
          order = dayOrder.indexOf(day).takeIf { it >= 0 } ?: Int.MAX_VALUE,
          periods = linkedSetOf(),
        )
      }

      group.periods.add(periodValue)
    }

    val formatted = mutableListOf<String>()

    groups.values
      .sortedWith(
        compareBy<TimeTableGroup> { it.order }
          .thenBy {
            // roomLabel empty first
            if (it.roomLabel.isBlank()) "" else it.roomLabel
          }
          .thenBy { it.roomLabel },
      )
      .forEach { group ->
        val periods = group.periods.toList().sorted()
        compressPeriods(periods).forEach { range ->
          val periodLabel = if (range.start == range.end) {
            "${range.start}"
          } else {
            "${range.start}-${range.end}"
          }

          val roomSuffix = if (includeRoom && group.roomLabel.isNotBlank()) "(${group.roomLabel})" else ""
          formatted.add("${group.dayLabel}$periodLabel$roomSuffix")
        }
      }

    val result = (formatted + fallbacks)
      .filter { it.isNotBlank() }
      .distinct()

    return result.joinToString(", ")
  }

  private data class TimeTableGroup(
    val day: DayOfWeek,
    val dayLabel: String,
    val roomLabel: String,
    val order: Int,
    val periods: LinkedHashSet<Int>,
  )

  private data class PeriodRange(
    val start: Int,
    val end: Int,
  )

  private fun compressPeriods(periods: List<Int>): List<PeriodRange> {
    if (periods.isEmpty()) return emptyList()

    val ranges = mutableListOf<PeriodRange>()
    var start = periods.first()
    var end = start

    for (index in 1 until periods.size) {
      val current = periods[index]
      if (current == end + 1) {
        end = current
        continue
      }

      ranges.add(PeriodRange(start = start, end = end))
      start = current
      end = current
    }

    ranges.add(PeriodRange(start = start, end = end))
    return ranges
  }

  private fun buildFallbackLabel(
    dayLabel: String,
    period: Int?,
    roomName: String?,
    includeRoom: Boolean,
  ): String {
    val dayPart = dayLabel
    val periodPart = if (period != null && period > 0) period.toString() else ""
    val roomPart = if (includeRoom && !roomName.isNullOrBlank()) "(${roomName})" else ""

    if (dayPart.isBlank() && periodPart.isBlank() && roomPart.isBlank()) return ""
    return "$dayPart$periodPart$roomPart"
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
}
