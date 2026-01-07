package com.github.kavos113.desy.ui.list

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.Modifier
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Color
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
  onOpenSearch: () -> Unit,
  onSelectLecture: (Int) -> Unit,
  modifier: Modifier = Modifier,
) {
  Column(
    modifier = modifier
      .fillMaxSize()
      .padding(12.dp),
  ) {
    Row(modifier = Modifier.fillMaxWidth()) {
      Spacer(modifier = Modifier.weight(1f))
      Button(onClick = onOpenSearch) {
        Text("検索")
      }
    }

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
          verticalArrangement = Arrangement.spacedBy(0.dp),
        ) {
          items(uiState.items, key = { it.id }) { item ->
            LectureListRow(item, onClick = { onSelectLecture(item.id) })
          }
        }
      }
    }
  }
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
    LectureListScreen(uiState = sample, onOpenSearch = {}, onSelectLecture = {})
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListHeaderPreview() {
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LectureListHeader()
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListRowPreview() {
  val item = LectureSummary(
    id = 1,
    university = "Sample University",
    title = "線形代数",
    department = "理学院",
    timetables = listOf(
      TimeTable(
        lectureId = 1,
        semester = Semester.spring,
        dayOfWeek = DayOfWeek.tuesday,
        period = 2,
      ),
    ),
  )

  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LectureListHeader()
      LectureListRow(item, onClick = {})
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListEmptyPreview() {
  DesyTheme {
    LectureListScreen(
      uiState = LectureListUiState(isLoading = false, items = emptyList()),
      onOpenSearch = {},
      onSelectLecture = {},
    )
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListLoadingPreview() {
  DesyTheme {
    LectureListScreen(
      uiState = LectureListUiState(isLoading = true),
      onOpenSearch = {},
      onSelectLecture = {},
    )
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureListErrorPreview() {
  DesyTheme {
    LectureListScreen(
      uiState = LectureListUiState(isLoading = false, errorMessage = "エラーが発生しました"),
      onOpenSearch = {},
      onSelectLecture = {},
    )
  }
}

@Composable
private fun LectureListHeader() {
  LectureListTableRow(
    cells = listOf(
      LectureListCellSpec(text = "講義名", weight = 0.5f),
      LectureListCellSpec(text = "時間割", weight = 0.2f),
      LectureListCellSpec(text = "開講期", weight = 0.1f),
      LectureListCellSpec(text = "学部/学科", weight = 0.2f),
    ),
    isHeader = true,
    onClick = null,
  )
}

@Composable
private fun LectureListRow(item: LectureSummary, onClick: () -> Unit) {
  LectureListTableRow(
    cells = listOf(
      LectureListCellSpec(text = item.title, weight = 0.5f),
      LectureListCellSpec(text = formatLectureTimetable(item.timetables), weight = 0.2f),
      LectureListCellSpec(text = formatLectureOpenTerm(item.timetables), weight = 0.1f),
      LectureListCellSpec(text = item.department.orEmpty(), weight = 0.2f),
    ),
    isHeader = false,
    onClick = onClick,
  )
}

private data class LectureListCellSpec(
  val text: String,
  val weight: Float,
)

@Composable
private fun LectureListTableRow(
  cells: List<LectureListCellSpec>,
  isHeader: Boolean,
  onClick: (() -> Unit)?,
  modifier: Modifier = Modifier,
) {
  val outline = MaterialTheme.colorScheme.outline
  Row(
    modifier = modifier
      .fillMaxWidth()
      .then(if (onClick != null) Modifier.clickable(onClick = onClick) else Modifier),
    verticalAlignment = Alignment.CenterVertically,
  ) {
    cells.forEachIndexed { index, cell ->
      TableCell(
        text = cell.text,
        weight = cell.weight,
        outlineColor = outline,
        isFirstColumn = index == 0,
        drawTop = isHeader,
        isHeader = isHeader,
      )
    }
  }
}

@Composable
private fun RowScope.TableCell(
  text: String,
  weight: Float,
  outlineColor: Color,
  isFirstColumn: Boolean,
  drawTop: Boolean,
  isHeader: Boolean,
) {
  val stroke = 1.dp
  Box(
    modifier = Modifier
      .weight(weight)
      .drawBehind {
        val s = stroke.toPx()
        // top border only for the first(header) row
        if (drawTop) {
          drawLine(outlineColor, start = Offset(0f, 0f), end = Offset(size.width, 0f), strokeWidth = s)
        }
        // bottom border for every row
        drawLine(outlineColor, start = Offset(0f, size.height), end = Offset(size.width, size.height), strokeWidth = s)
        // left border only for first column
        if (isFirstColumn) {
          drawLine(outlineColor, start = Offset(0f, 0f), end = Offset(0f, size.height), strokeWidth = s)
        }
        // right border for every cell
        drawLine(outlineColor, start = Offset(size.width, 0f), end = Offset(size.width, size.height), strokeWidth = s)
      }
      .padding(horizontal = 3.dp, vertical = 4.dp),
  ) {
    Text(
      text = text,
      style = if (isHeader) MaterialTheme.typography.labelLarge else MaterialTheme.typography.bodyMedium,
      maxLines = 1,
      overflow = TextOverflow.Clip,
    )
  }
}


