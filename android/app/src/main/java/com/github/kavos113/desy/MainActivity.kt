package com.github.kavos113.desy

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import com.github.kavos113.desy.ui.detail.LectureDetailScreen
import com.github.kavos113.desy.ui.list.LectureListScreen
import com.github.kavos113.desy.ui.search.LectureSearchScreen
import com.github.kavos113.desy.ui.theme.DesyTheme
import com.github.kavos113.desy.ui.viewmodel.LectureViewModel
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
  private val lectureViewModel: LectureViewModel by viewModels()

  override fun onCreate(savedInstanceState: Bundle?) {
    super.onCreate(savedInstanceState)
    enableEdgeToEdge()
    setContent {
      DesyTheme {
        var isSearchOpen by rememberSaveable { mutableStateOf(false) }
        var selectedLectureId by rememberSaveable { mutableStateOf<Int?>(null) }

        Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
          if (selectedLectureId != null) {
            LectureDetailScreen(
              uiState = lectureViewModel.detailUiState,
              onBack = {
                selectedLectureId = null
              },
              onSelectRelatedLecture = { lectureId ->
                selectedLectureId = lectureId
                lectureViewModel.loadLectureDetails(lectureId)
              },
              modifier = Modifier.padding(innerPadding),
            )
          } else if (isSearchOpen) {
            LectureSearchScreen(
              onSearch = { query ->
                lectureViewModel.loadLectures(query)
                isSearchOpen = false
              },
              onCancel = {
                isSearchOpen = false
              },
              modifier = Modifier.padding(innerPadding),
            )
          } else {
            LectureListScreen(
              uiState = lectureViewModel.uiState,
              onOpenSearch = { isSearchOpen = true },
              onSelectLecture = { lectureId ->
                selectedLectureId = lectureId
                lectureViewModel.loadLectureDetails(lectureId)
              },
              modifier = Modifier.padding(innerPadding),
            )
          }
        }
      }
    }
  }
}
