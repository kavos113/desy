package com.github.kavos113.desy.repository.db.dao

import androidx.room.Dao
import androidx.room.Query
import androidx.room.RawQuery
import androidx.sqlite.db.SupportSQLiteQuery
import com.github.kavos113.desy.repository.db.entity.LectureEntity
import com.github.kavos113.desy.repository.db.entity.LecturePlanEntity
import com.github.kavos113.desy.repository.db.model.LectureSummaryRow
import com.github.kavos113.desy.repository.db.model.TeacherRow
import com.github.kavos113.desy.repository.db.model.TimetableRow

@Dao
interface LectureDao {
  @Query("SELECT * FROM lectures WHERE id = :id")
  suspend fun getLectureById(id: Int): LectureEntity?

  @RawQuery
  suspend fun searchLectureSummaries(query: SupportSQLiteQuery): List<LectureSummaryRow>

  @Query(
    """
      SELECT tt.lecture_id, tt.semester, tt.room_id, r.name AS room_name, tt.day_of_week, tt.period
      FROM timetables tt
      LEFT JOIN rooms r ON r.id = tt.room_id
      WHERE tt.lecture_id IN (:lectureIds)
      ORDER BY tt.lecture_id, tt.semester, tt.day_of_week, tt.period
    """
  )
  suspend fun getTimetablesForLectures(lectureIds: List<Int>): List<TimetableRow>

  @Query(
    """
      SELECT tt.lecture_id, tt.semester, tt.room_id, r.name AS room_name, tt.day_of_week, tt.period
      FROM timetables tt
      LEFT JOIN rooms r ON r.id = tt.room_id
      WHERE tt.lecture_id = :lectureId
      ORDER BY tt.semester, tt.day_of_week, tt.period
    """
  )
  suspend fun getTimetablesForLecture(lectureId: Int): List<TimetableRow>

  @Query(
    """
      SELECT lt.lecture_id, t.id, t.name, t.url
      FROM lecture_teachers lt
      JOIN teachers t ON t.id = lt.teacher_id
      WHERE lt.lecture_id IN (:lectureIds)
      ORDER BY lt.lecture_id, t.id
    """
  )
  suspend fun getTeachersForLectures(lectureIds: List<Int>): List<TeacherRow>

  @Query(
    """
      SELECT lt.lecture_id, t.id, t.name, t.url
      FROM lecture_teachers lt
      JOIN teachers t ON t.id = lt.teacher_id
      WHERE lt.lecture_id = :lectureId
      ORDER BY t.id
    """
  )
  suspend fun getTeachersForLecture(lectureId: Int): List<TeacherRow>

  @Query("SELECT * FROM lecture_plans WHERE lecture_id = :lectureId ORDER BY count")
  suspend fun getLecturePlans(lectureId: Int): List<LecturePlanEntity>

  @Query("SELECT keyword FROM lecture_keywords WHERE lecture_id = :lectureId ORDER BY keyword")
  suspend fun getKeywords(lectureId: Int): List<String>

  @Query("SELECT related_lecture_id FROM related_courses WHERE lecture_id = :lectureId ORDER BY related_lecture_id")
  suspend fun getRelatedLectureIds(lectureId: Int): List<Int>

  @Query("SELECT code FROM related_course_codes WHERE lecture_id = :lectureId ORDER BY code")
  suspend fun getRelatedCourseCodes(lectureId: Int): List<String>
}
