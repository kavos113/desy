package com.github.kavos113.desy.usecase

import com.github.kavos113.desy.domain.Lecture
import com.github.kavos113.desy.domain.LectureRepository
import com.github.kavos113.desy.domain.LectureSummary
import com.github.kavos113.desy.domain.SearchQuery

interface LectureUsecase {
  suspend fun searchLectures(query: SearchQuery): List<LectureSummary>
  suspend fun getLectureDetails(lectureId: Int): Lecture?
}

class DefaultLectureUsecase(
  private val lectureRepository: LectureRepository,
) : LectureUsecase {
  override suspend fun searchLectures(query: SearchQuery): List<LectureSummary> {
    return lectureRepository.search(query)
  }

  override suspend fun getLectureDetails(lectureId: Int): Lecture? {
    if (lectureId <= 0) return null
    return lectureRepository.findById(lectureId)
  }
}
