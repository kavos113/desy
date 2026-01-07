package com.github.kavos113.desy.repository

import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Lecture
import com.github.kavos113.desy.domain.LecturePlan
import com.github.kavos113.desy.domain.LectureRepository
import com.github.kavos113.desy.domain.LectureSummary
import com.github.kavos113.desy.domain.LectureType
import com.github.kavos113.desy.domain.Level
import com.github.kavos113.desy.domain.Room
import com.github.kavos113.desy.domain.SearchQuery
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.Teacher
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.repository.db.DesyDatabase
import javax.inject.Inject

class RoomLectureRepository @Inject constructor(
  private val db: DesyDatabase,
) : LectureRepository {
  private val lectureDao = db.lectureDao()

  override suspend fun findById(id: Int): Lecture? {
    val lecture = lectureDao.getLectureById(id) ?: return null

    val timetables = lectureDao.getTimetablesForLecture(id).map { row ->
      TimeTable(
        lectureId = row.lectureId,
        semester = Semester.fromDb(row.semester),
        room = if (row.roomId == null && row.roomName == null) null else Room(id = row.roomId, name = row.roomName),
        dayOfWeek = DayOfWeek.fromDb(row.dayOfWeek),
        period = row.period,
      )
    }

    val teachers = lectureDao.getTeachersForLecture(id).map { row ->
      Teacher(
        id = row.teacherId,
        name = row.name,
        url = row.url,
      )
    }

    val plans = lectureDao.getLecturePlans(id).map { p ->
      LecturePlan(count = p.count, plan = p.plan, assignment = p.assignment)
    }

    val keywords = lectureDao.getKeywords(id)
    val related = lectureDao.getRelatedLectureIds(id)
    val relatedCodes = lectureDao.getRelatedCourseCodes(id)

    return Lecture(
      id = lecture.id,
      university = lecture.university,
      title = lecture.title,
      englishTitle = lecture.englishTitle,
      department = lecture.department,
      lectureType = LectureType.fromDb(lecture.lectureType),
      code = lecture.code,
      level = Level.fromDb(lecture.level),
      credit = lecture.credit,
      year = lecture.year,
      openTerm = lecture.openTerm,
      language = lecture.language,
      url = lecture.url,
      abstractText = lecture.abstractText,
      goal = lecture.goal,
      experience = lecture.experience,
      flow = lecture.flow,
      outOfClassWork = lecture.outOfClassWork,
      textbook = lecture.textbook,
      referenceBook = lecture.referenceBook,
      assessment = lecture.assessment,
      prerequisite = lecture.prerequisite,
      contact = lecture.contact,
      officeHours = lecture.officeHours,
      note = lecture.note,
      updatedAt = lecture.updatedAt,
      timetables = timetables,
      teachers = teachers,
      lecturePlans = plans,
      keywords = keywords,
      relatedCourses = related,
      relatedCourseCodes = relatedCodes,
    )
  }

  override suspend fun search(query: SearchQuery): List<LectureSummary> {
    val built = LectureSearchQueryBuilder.build(query)
    val rows = lectureDao.searchLectureSummaries(built.toSupportQuery())
    if (rows.isEmpty()) return emptyList()

    val ids = rows.map { it.id }

    val timetablesMap = lectureDao.getTimetablesForLectures(ids)
      .groupBy { it.lectureId }
      .mapValues { (_, list) ->
        list.map { row ->
          TimeTable(
            lectureId = row.lectureId,
            semester = Semester.fromDb(row.semester),
            room = if (row.roomId == null && row.roomName == null) null else Room(id = row.roomId, name = row.roomName),
            dayOfWeek = DayOfWeek.fromDb(row.dayOfWeek),
            period = row.period,
          )
        }
      }

    val teachersMap = lectureDao.getTeachersForLectures(ids)
      .groupBy { it.lectureId }
      .mapValues { (_, list) ->
        list.map { row ->
          Teacher(id = row.teacherId, name = row.name, url = row.url)
        }
      }

    val summaries = rows.map { row ->
      LectureSummary(
        id = row.id,
        university = row.university,
        title = row.title,
        department = row.department.takeIf { it.isNotBlank() },
        code = row.code.takeIf { it.isNotBlank() },
        level = Level.fromDb(row.level),
        credit = row.credit,
        year = row.year,
        timetables = timetablesMap[row.id].orEmpty(),
        teachers = teachersMap[row.id].orEmpty(),
      )
    }

    return if (query.filterNotResearch) {
      NotResearchFilter.filter(summaries)
    } else {
      summaries
    }
  }
}
