package com.github.kavos113.desy.domain

/**
 * バックエンド(domain)のモデルに合わせた、Android側のドメインモデル。
 * Android版はスクレイピングを行わず、DB(SQLite)を読み込んで表示する用途のため、
 * 取得に必要な読み取り系のみを定義しています。
 */

data class Lecture(
  val id: Int,
  val university: String,
  val title: String,
  val englishTitle: String? = null,
  val department: String? = null,
  val lectureType: LectureType? = null,
  val code: String? = null,
  val level: Level? = null,
  val credit: Int? = null,
  val year: Int? = null,
  val openTerm: String? = null,
  val language: String? = null,
  val url: String? = null,
  val abstractText: String? = null,
  val goal: String? = null,
  val experience: String? = null,
  val flow: String? = null,
  val outOfClassWork: String? = null,
  val textbook: String? = null,
  val referenceBook: String? = null,
  val assessment: String? = null,
  val prerequisite: String? = null,
  val contact: String? = null,
  val officeHours: String? = null,
  val note: String? = null,
  val updatedAt: String? = null,
  val timetables: List<TimeTable> = emptyList(),
  val teachers: List<Teacher> = emptyList(),
  val lecturePlans: List<LecturePlan> = emptyList(),
  val keywords: List<String> = emptyList(),
  val relatedCourseCodes: List<String> = emptyList(),
  val relatedCourses: List<Int> = emptyList(),
)
