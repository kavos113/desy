package com.github.kavos113.desy.repository

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.SearchQuery
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.TimeTable
import org.junit.Assert.assertTrue
import org.junit.Test

class LectureSearchQueryBuilderTest {
  @Test
  fun build_includesExpectedJoinsAndConditions() {
    val query = SearchQuery(
      title = "アルゴリズム",
      teacherName = "山田",
      keywords = listOf("データ", "構造"),
      room = "E",
      semesters = listOf(Semester.spring, Semester.fall),
      timetables = listOf(
        TimeTable(lectureId = 0, dayOfWeek = DayOfWeek.monday, period = 2),
        TimeTable(lectureId = 0, dayOfWeek = DayOfWeek.friday, period = 5),
      ),
    )

    val built = LectureSearchQueryBuilder.build(query)

    assertTrue(built.sql.contains("FROM lectures l"))
    assertTrue(built.sql.contains("JOIN lecture_teachers"))
    assertTrue(built.sql.contains("JOIN teachers"))
    assertTrue(built.sql.contains("JOIN lecture_keywords"))
    assertTrue(built.sql.contains("JOIN timetables tt"))
    assertTrue(built.sql.contains("JOIN rooms r"))
    assertTrue(built.sql.contains("t.name LIKE ?"))
    assertTrue(built.sql.contains("lk.keyword IN"))
    assertTrue(built.sql.contains("tt.semester IN"))
    assertTrue(built.sql.contains("r.name LIKE ?"))
    assertTrue(built.sql.contains("l.title LIKE ?"))
    assertTrue(built.sql.contains("ORDER BY l.year DESC"))

    // args数のざっくり検証（プレースホルダ数と一致するはず）
    val placeholderCount = built.sql.count { it == '?' }
    assertTrue(placeholderCount == built.args.size)
  }
}
