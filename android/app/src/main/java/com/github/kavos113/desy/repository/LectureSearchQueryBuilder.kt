package com.github.kavos113.desy.repository

import androidx.sqlite.db.SimpleSQLiteQuery
import androidx.sqlite.db.SupportSQLiteQuery
import com.github.kavos113.desy.domain.SearchQuery

/**
 * backend/presentation/repository/sqlite/lecture_repository.go のSearch組み立てに合わせたSQLビルダー。
 */
object LectureSearchQueryBuilder {
  data class BuiltQuery(
    val sql: String,
    val args: List<Any?>,
  ) {
    fun toSupportQuery(): SupportSQLiteQuery = SimpleSQLiteQuery(sql, args.toTypedArray())
  }

  fun build(query: SearchQuery): BuiltQuery {
    val selectBuilder = StringBuilder()
    // NULL文字列が混ざるとRoom側DTOで扱いづらいのでIFNULLで吸収
    selectBuilder.append(
      "SELECT DISTINCT l.id, l.university, l.title, IFNULL(l.department, '') AS department, IFNULL(l.code, '') AS code, l.level, l.credit, l.year FROM lectures l"
    )

    val joins = mutableListOf<String>()
    val conditions = mutableListOf<String>()
    val args = mutableListOf<Any?>()

    if (query.teacherName.isNotBlank()) {
      joins += "JOIN lecture_teachers lt ON lt.lecture_id = l.id JOIN teachers t ON t.id = lt.teacher_id"
      conditions += "t.name LIKE ?"
      args.add("%${query.teacherName}%")
    }

    if (query.keywords.isNotEmpty()) {
      joins += "JOIN lecture_keywords lk ON lk.lecture_id = l.id"
      conditions += "lk.keyword IN (${placeholders(query.keywords.size)})"
      args.addAll(query.keywords)
    }

    val timetableJoinRequired = query.timetables.isNotEmpty() || query.room.isNotBlank() || query.semesters.isNotEmpty()
    if (timetableJoinRequired) {
      joins += "JOIN timetables tt ON tt.lecture_id = l.id"
    }

    if (query.timetables.isNotEmpty()) {
      val timetableFilters = mutableListOf<String>()
      for (tt in query.timetables) {
        val day = tt.dayOfWeek?.name
        val period = tt.period
        if (day.isNullOrBlank() && (period == null || period == 0)) continue
        timetableFilters += "(tt.day_of_week = ? AND tt.period = ?)"
        args.add(day)
        args.add(period)
      }
      if (timetableFilters.isNotEmpty()) {
        conditions += "(${timetableFilters.joinToString(" OR ")})"
      }
    }

    if (query.semesters.isNotEmpty()) {
      conditions += "tt.semester IN (${placeholders(query.semesters.size)})"
      args.addAll(query.semesters.map { it.name })
    }

    if (query.room.isNotBlank()) {
      joins += "JOIN rooms r ON r.id = tt.room_id"
      conditions += "r.name LIKE ?"
      args.add("%${query.room}%")
    }

    if (query.title.isNotBlank()) {
      conditions += "(l.title LIKE ? OR IFNULL(l.english_title, '') LIKE ?)"
      val like = "%${query.title}%"
      args.add(like)
      args.add(like)
    }

    if (query.departments.isNotEmpty()) {
      conditions += "l.department IN (${placeholders(query.departments.size)})"
      args.addAll(query.departments)
    }

    if (query.year != 0) {
      conditions += "l.year = ?"
      args.add(query.year)
    }

    if (query.levels.isNotEmpty()) {
      conditions += "l.level IN (${placeholders(query.levels.size)})"
      args.addAll(query.levels.map { it.value })
    }

    if (joins.isNotEmpty()) {
      selectBuilder.append(' ')
      selectBuilder.append(joins.joinToString(" "))
    }

    if (conditions.isNotEmpty()) {
      selectBuilder.append(" WHERE ")
      selectBuilder.append(conditions.joinToString(" AND "))
    }

    selectBuilder.append(" ORDER BY l.year DESC, l.title ASC")

    return BuiltQuery(selectBuilder.toString(), args)
  }

  private fun placeholders(count: Int): String = List(count) { "?" }.joinToString(",")
}
