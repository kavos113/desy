package com.github.kavos113.desy

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.ui.Modifier
import com.github.kavos113.desy.ui.list.LectureListScreen
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
        Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
          LectureListScreen(
            uiState = lectureViewModel.uiState,
            modifier = Modifier.padding(innerPadding),
          )
        }
      }
    }
  }
}
