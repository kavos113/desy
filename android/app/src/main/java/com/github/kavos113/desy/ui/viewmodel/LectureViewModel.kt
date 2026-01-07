package com.github.kavos113.desy.ui.viewmodel

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.github.kavos113.desy.domain.SearchQuery
import com.github.kavos113.desy.ui.detail.RelatedCourseEntry
import com.github.kavos113.desy.usecase.LectureUsecase
import dagger.hilt.android.lifecycle.HiltViewModel
import javax.inject.Inject
import kotlinx.coroutines.launch

/**
 * 講義一覧/検索画面で共通利用するViewModel。
 * 現時点では検索UIは未実装のため、空条件(SearchQuery())で一覧を取得します。
 */
@HiltViewModel
class LectureViewModel @Inject constructor(
  private val lectureUsecase: LectureUsecase,
) : ViewModel() {

  var uiState: LectureListUiState by mutableStateOf(LectureListUiState())
    private set

  var detailUiState: LectureDetailUiState by mutableStateOf(LectureDetailUiState())
    private set

  fun loadLectures(query: SearchQuery) {
    uiState = uiState.copy(query = query, isLoading = true, errorMessage = null)

    viewModelScope.launch {
      runCatching {
        lectureUsecase.searchLectures(query)
      }.onSuccess { items ->
        uiState = uiState.copy(items = items, isLoading = false)
      }.onFailure { throwable ->
        uiState = uiState.copy(
          items = emptyList(),
          isLoading = false,
          errorMessage = throwable.message ?: "講義一覧の取得に失敗しました。",
        )
      }
    }
  }

  fun loadLectureDetails(lectureId: Int) {
    detailUiState = detailUiState.copy(
      lectureId = lectureId,
      lecture = null,
      relatedCourses = emptyList(),
      isLoading = true,
      errorMessage = null,
    )

    viewModelScope.launch {
      runCatching {
        lectureUsecase.getLectureDetails(lectureId)
      }.onSuccess { lecture ->
        if (lecture == null) {
          detailUiState = detailUiState.copy(
            lecture = null,
            relatedCourses = emptyList(),
            isLoading = false,
            errorMessage = "講義が見つかりませんでした。",
          )
        } else {
          val relatedEntries = buildList {
            lecture.relatedCourseCodes
              .map { it.trim() }
              .filter { it.isNotEmpty() }
              .forEach { code -> add(RelatedCourseEntry(code = code)) }

            lecture.relatedCourses
              .distinct()
              .filter { it > 0 }
              .forEach { relatedId ->
                val relatedLecture = runCatching {
                  lectureUsecase.getLectureDetails(relatedId)
                }.getOrNull()

                val code = relatedLecture?.code?.trim().orEmpty().ifBlank { relatedId.toString() }
                val title = relatedLecture?.title
                add(RelatedCourseEntry(code = code, id = relatedId, title = title))
              }
          }

          detailUiState = detailUiState.copy(
            lecture = lecture,
            relatedCourses = relatedEntries,
            isLoading = false,
          )
        }
      }.onFailure { throwable ->
        detailUiState = detailUiState.copy(
          lecture = null,
          relatedCourses = emptyList(),
          isLoading = false,
          errorMessage = throwable.message ?: "講義詳細の取得に失敗しました。",
        )
      }
    }
  }
}
