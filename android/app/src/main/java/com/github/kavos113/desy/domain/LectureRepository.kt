package com.github.kavos113.desy.domain

interface LectureRepository {
  suspend fun findById(id: Int): Lecture?
  suspend fun search(query: SearchQuery): List<LectureSummary>
}
