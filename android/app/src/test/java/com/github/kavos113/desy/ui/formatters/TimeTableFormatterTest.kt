package com.github.kavos113.desy.ui.formatters

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Room
import com.github.kavos113.desy.domain.TimeTable
import org.junit.Assert.assertEquals
import org.junit.Test

class TimeTableFormatterTest {
  @Test
  fun format_mergesConsecutivePeriods_sameDaySameRoom() {
    val result = TimeTableFormatter.format(
      timetables = listOf(
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.monday, period = 1, room = Room(name = "W1")),
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.monday, period = 2, room = Room(name = "W1")),
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.monday, period = 3, room = Room(name = "W1")),
      ),
    )

    assertEquals("月1-3(W1)", result)
  }

  @Test
  fun format_canOmitRoom() {
    val result = TimeTableFormatter.format(
      timetables = listOf(
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.friday, period = 5, room = Room(name = "M110")),
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.friday, period = 6, room = Room(name = "M110")),
      ),
      includeRoom = false,
    )

    assertEquals("金5-6", result)
  }

  @Test
  fun format_separatesDifferentRooms_evenSameDay() {
    val result = TimeTableFormatter.format(
      timetables = listOf(
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.tuesday, period = 3, room = Room(name = "W1")),
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.tuesday, period = 4, room = Room(name = "W2")),
      ),
    )

    assertEquals("火3(W1), 火4(W2)", result)
  }

  @Test
  fun format_doesNotMergeNonConsecutivePeriods() {
    val result = TimeTableFormatter.format(
      timetables = listOf(
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.wednesday, period = 2, room = Room(name = "W1")),
        TimeTable(lectureId = 1, dayOfWeek = DayOfWeek.wednesday, period = 4, room = Room(name = "W1")),
      ),
    )

    assertEquals("水2(W1), 水4(W1)", result)
  }
}
