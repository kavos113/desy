package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "lectures")
data class LectureEntity(
  @PrimaryKey
  @ColumnInfo(name = "id")
  val id: Int,

  @ColumnInfo(name = "university")
  val university: String,

  @ColumnInfo(name = "title")
  val title: String,

  @ColumnInfo(name = "english_title")
  val englishTitle: String?,

  @ColumnInfo(name = "department")
  val department: String?,

  @ColumnInfo(name = "lecture_type")
  val lectureType: String?,

  @ColumnInfo(name = "code")
  val code: String?,

  @ColumnInfo(name = "level")
  val level: Int?,

  @ColumnInfo(name = "credit")
  val credit: Int?,

  @ColumnInfo(name = "year")
  val year: Int?,

  @ColumnInfo(name = "open_term")
  val openTerm: String?,

  @ColumnInfo(name = "language")
  val language: String?,

  @ColumnInfo(name = "url")
  val url: String?,

  @ColumnInfo(name = "abstract")
  val abstractText: String?,

  @ColumnInfo(name = "goal")
  val goal: String?,

  @ColumnInfo(name = "experience")
  val experience: String?,

  @ColumnInfo(name = "flow")
  val flow: String?,

  @ColumnInfo(name = "out_of_class_work")
  val outOfClassWork: String?,

  @ColumnInfo(name = "textbook")
  val textbook: String?,

  @ColumnInfo(name = "reference_book")
  val referenceBook: String?,

  @ColumnInfo(name = "assessment")
  val assessment: String?,

  @ColumnInfo(name = "prerequisite")
  val prerequisite: String?,

  @ColumnInfo(name = "contact")
  val contact: String?,

  @ColumnInfo(name = "office_hours")
  val officeHours: String?,

  @ColumnInfo(name = "note")
  val note: String?,

  @ColumnInfo(name = "updated_at")
  val updatedAt: String?,
)
