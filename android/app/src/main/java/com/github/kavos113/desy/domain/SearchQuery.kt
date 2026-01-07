package com.github.kavos113.desy.domain

data class SearchQuery(
	val title: String = "",
	val keywords: List<String> = emptyList(),
	val departments: List<String> = emptyList(),
	val year: Int = 0,
	val teacherName: String = "",
	val room: String = "",
	val semesters: List<Semester> = emptyList(),
	val timetables: List<TimeTable> = emptyList(),
	val levels: List<Level> = emptyList(),
	val filterNotResearch: Boolean = false,
)
