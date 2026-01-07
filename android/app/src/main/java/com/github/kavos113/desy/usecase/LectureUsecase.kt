package com.github.kavos113.desy.usecase

import com.github.kavos113.desy.domain.Lecture
import com.github.kavos113.desy.domain.LectureRepository
import com.github.kavos113.desy.domain.LectureSummary
import com.github.kavos113.desy.domain.SearchQuery
import javax.inject.Inject

class LectureUsecase @Inject constructor(
  private val lectureRepository: LectureRepository,
) {
  suspend fun searchLectures(query: SearchQuery): List<LectureSummary> {
    return lectureRepository.search(query)
  }

  suspend fun getLectureDetails(lectureId: Int): Lecture? {
    if (lectureId <= 0) return null
    return lectureRepository.findById(lectureId)
  }
}
