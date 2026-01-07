package com.github.kavos113.desy.ui.viewmodel

import com.github.kavos113.desy.domain.LectureSummary
import com.github.kavos113.desy.domain.SearchQuery

data class LectureListUiState(
  val query: SearchQuery = SearchQuery(),
  val items: List<LectureSummary> = emptyList(),
  val isLoading: Boolean = false,
  val errorMessage: String? = null,
)
