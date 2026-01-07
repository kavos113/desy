package com.github.kavos113.desy.ui.list

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.LectureSummary
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.ui.viewmodel.LectureListUiState
import com.github.kavos113.desy.ui.theme.DesyTheme

@Composable
fun LectureListScreen(
  uiState: LectureListUiState,
  modifier: Modifier = Modifier,
) {
  Column(
    modifier = modifier
      .fillMaxSize()
      .padding(12.dp),
    verticalArrangement = Arrangement.spacedBy(8.dp),
  ) {
    LectureListHeader()

    when {
      uiState.isLoading -> {
        Text("読み込み中…")
      }

      uiState.errorMessage != null -> {
        Text(uiState.errorMessage)
      }

      uiState.items.isEmpty() -> {
        Text("講義が見つかりませんでした。")
      }

      else -> {
        LazyColumn(
          modifier = Modifier.fillMaxWidth(),
          verticalArrangement = Arrangement.spacedBy(6.dp),
        ) {
          items(uiState.items, key = { it.id }) { item ->
            LectureListRow(item)
          }
        }
      }
    }
  }
}

@Composable
private fun LectureListHeader() {
  Row(
    modifier = Modifier
      .fillMaxWidth()
      .padding(vertical = 6.dp),
    horizontalArrangement = Arrangement.spacedBy(8.dp),
  ) {
    HeaderCell("講義名", weight = 0.42f)
    HeaderCell("時間割", weight = 0.22f)
    HeaderCell("開講期", weight = 0.18f)
    HeaderCell("学部/学科", weight = 0.18f)
  }
}

@Composable
private fun LectureListRow(item: LectureSummary) {
  Row(
    modifier = Modifier.fillMaxWidth(),
    horizontalArrangement = Arrangement.spacedBy(8.dp),
  ) {
    BodyCell(item.title, weight = 0.42f)
    BodyCell(formatLectureTimetable(item.timetables), weight = 0.22f)
    BodyCell(formatLectureOpenTerm(item.timetables), weight = 0.18f)
    BodyCell(item.department.orEmpty(), weight = 0.18f)
  }
}

@Composable
private fun RowScope.HeaderCell(text: String, weight: Float) {
  Text(
    text = text,
    style = MaterialTheme.typography.labelLarge,
    modifier = Modifier.weight(weight),
    maxLines = 1,
    overflow = TextOverflow.Ellipsis,
  )
}

@Composable
private fun RowScope.BodyCell(text: String, weight: Float) {
  Text(
    text = text,
    style = MaterialTheme.typography.bodyMedium,
    modifier = Modifier.weight(weight),
    maxLines = 2,
    overflow = TextOverflow.Ellipsis,
  )
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListScreenPreview() {
  val sample = LectureListUiState(
    isLoading = false,
    items = listOf(
      LectureSummary(
        id = 1,
        university = "Sample University",
        title = "プログラミング基礎",
        department = "情報学部",
        timetables = listOf(
          TimeTable(
            lectureId = 1,
            semester = Semester.spring,
            dayOfWeek = DayOfWeek.monday,
            period = 1,
          ),
          TimeTable(
            lectureId = 1,
            semester = Semester.fall,
            dayOfWeek = DayOfWeek.wednesday,
            period = 2,
          ),
        ),
      ),
      LectureSummary(
        id = 2,
        university = "Sample University",
        title = "データベース",
        department = "工学部",
        timetables = listOf(
          TimeTable(
            lectureId = 2,
            semester = Semester.summer,
            dayOfWeek = DayOfWeek.thursday,
            period = 3,
          ),
        ),
      ),
    ),
  )

  DesyTheme {
    LectureListScreen(uiState = sample)
  }
}
