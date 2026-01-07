package com.github.kavos113.desy.ui.viewmodel

import com.github.kavos113.desy.domain.Lecture
import com.github.kavos113.desy.ui.detail.RelatedCourseEntry

data class LectureDetailUiState(
  val lectureId: Int? = null,
  val lecture: Lecture? = null,
  val relatedCourses: List<RelatedCourseEntry> = emptyList(),
  val isLoading: Boolean = false,
  val errorMessage: String? = null,
)
