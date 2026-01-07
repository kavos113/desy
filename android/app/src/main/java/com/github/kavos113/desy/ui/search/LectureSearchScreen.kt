package com.github.kavos113.desy.ui.search

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Checkbox
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Level
import com.github.kavos113.desy.domain.SearchQuery
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.ui.theme.DesyTheme

@Composable
fun LectureSearchScreen(
  onSearch: (SearchQuery) -> Unit,
  onCancel: () -> Unit,
  modifier: Modifier = Modifier,
) {
  LectureSearchScreenContent(onSearch = onSearch, onCancel = onCancel, modifier = modifier)
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureSearchScreenPreview() {
  DesyTheme {
    LectureSearchScreen(onSearch = {}, onCancel = {})
  }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun LectureSearchScreenContent(
  onSearch: (SearchQuery) -> Unit,
  onCancel: () -> Unit,
  modifier: Modifier = Modifier,
) {
  var title by remember { mutableStateOf("") }
  var teacherName by remember { mutableStateOf("") }
  var room by remember { mutableStateOf("") }
  var yearText by remember { mutableStateOf("") }
  var departmentsText by remember { mutableStateOf("") }
  var keywordsText by remember { mutableStateOf("") }
  var filterNotResearch by remember { mutableStateOf(false) }

  val selectedLevels = remember { mutableStateListOf<Level>() }
  val selectedSemesters = remember { mutableStateListOf<Semester>() }
  val selectedTimetables = remember { mutableStateListOf<Pair<DayOfWeek, Int>>() }

  Column(
    modifier = modifier
      .fillMaxSize()
      .padding(12.dp),
    verticalArrangement = Arrangement.spacedBy(12.dp),
  ) {
    Row(modifier = Modifier.fillMaxWidth()) {
      Text("検索", style = MaterialTheme.typography.titleLarge)
      Spacer(modifier = Modifier.weight(1f))
      TextButton(onClick = onCancel) {
        Text("戻る")
      }
    }

    OutlinedTextField(
      value = title,
      onValueChange = { title = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("講義名") },
      singleLine = true,
    )

    OutlinedTextField(
      value = teacherName,
      onValueChange = { teacherName = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("教員名") },
      singleLine = true,
    )

    OutlinedTextField(
      value = room,
      onValueChange = { room = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("教室") },
      singleLine = true,
    )

    OutlinedTextField(
      value = yearText,
      onValueChange = { yearText = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("年度") },
      singleLine = true,
      placeholder = { Text("例: 2025") },
    )

    OutlinedTextField(
      value = departmentsText,
      onValueChange = { departmentsText = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("開講元（複数は空白/カンマ区切り）") },
    )

    OutlinedTextField(
      value = keywordsText,
      onValueChange = { keywordsText = it },
      modifier = Modifier.fillMaxWidth(),
      label = { Text("キーワード（複数は空白/カンマ区切り）") },
    )

    HorizontalDivider()

    Text("学年", style = MaterialTheme.typography.titleMedium)
    LevelCheckboxRow(selectedLevels)

    Spacer(Modifier.height(4.dp))

    Text("学期(Quarter)", style = MaterialTheme.typography.titleMedium)
    SemesterCheckboxRow(selectedSemesters)

    HorizontalDivider()

    Text("時間割", style = MaterialTheme.typography.titleMedium)
    TimetablePicker(
      selected = selectedTimetables,
    )

    HorizontalDivider()

    Row(
      modifier = Modifier.fillMaxWidth(),
      horizontalArrangement = Arrangement.spacedBy(12.dp),
    ) {
      Row(
        horizontalArrangement = Arrangement.spacedBy(8.dp),
      ) {
        Checkbox(checked = filterNotResearch, onCheckedChange = { filterNotResearch = it })
        Text("研究室配属系を除外")
      }

      Spacer(modifier = Modifier.weight(1f))

      Button(
        onClick = {
          val query = SearchQuery(
            title = title,
            teacherName = teacherName,
            room = room,
            year = parseYearInput(yearText),
            departments = parseDepartmentInput(departmentsText),
            keywords = parseKeywordInput(keywordsText),
            semesters = selectedSemesters.toList(),
            timetables = selectedTimetables.map { (day, period) ->
              TimeTable(lectureId = 0, dayOfWeek = day, period = period)
            },
            levels = selectedLevels.toList(),
            filterNotResearch = filterNotResearch,
          )
          onSearch(query)
        },
      ) {
        Text("Search")
      }
    }
  }
}

@Composable
private fun LevelCheckboxRow(selected: MutableList<Level>) {
  Column(verticalArrangement = Arrangement.spacedBy(4.dp)) {
    levelOptionRows().forEach { row ->
      Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
        row.forEach { (level, label) ->
          LabeledCheckbox(
            label = label,
            checked = selected.contains(level),
            onCheckedChange = { checked ->
              if (checked) selected.add(level) else selected.remove(level)
            },
          )
        }
      }
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LevelCheckboxRowPreview() {
  val selected = remember { mutableStateListOf(Level.bachelor1, Level.master1) }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      Text("学年", style = MaterialTheme.typography.titleMedium)
      LevelCheckboxRow(selected)
    }
  }
}

private fun levelOptionRows(): List<List<Pair<Level, String>>> {
  val all = listOf(
    Level.bachelor1 to "学士1年",
    Level.bachelor2 to "学士2年",
    Level.bachelor3 to "学士3年",
    Level.master1 to "修士1年",
    Level.master2 to "修士2年",
    Level.doctor to "博士課程",
  )
  return all.chunked(3)
}

@Composable
private fun SemesterCheckboxRow(selected: MutableList<Semester>) {
  Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
    listOf(
      Semester.spring to "1Q",
      Semester.summer to "2Q",
      Semester.fall to "3Q",
      Semester.winter to "4Q",
    ).forEach { (semester, label) ->
      LabeledCheckbox(
        label = label,
        checked = selected.contains(semester),
        onCheckedChange = { checked ->
          if (checked) selected.add(semester) else selected.remove(semester)
        },
      )
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun SemesterCheckboxRowPreview() {
  val selected = remember { mutableStateListOf(Semester.spring, Semester.fall) }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      Text("学期(Quarter)", style = MaterialTheme.typography.titleMedium)
      SemesterCheckboxRow(selected)
    }
  }
}

@Composable
private fun TimetablePicker(
  selected: MutableList<Pair<DayOfWeek, Int>>,
) {
  var dayExpanded by remember { mutableStateOf(false) }
  var periodExpanded by remember { mutableStateOf(false) }
  var selectedDay by remember { mutableStateOf(DayOfWeek.monday) }
  var selectedPeriod by remember { mutableStateOf(1) }

  Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
    Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
      Column {
        TextButton(onClick = { dayExpanded = true }) {
          Text("曜日: ${selectedDay.toJapaneseLabel()}")
        }
        DropdownMenu(expanded = dayExpanded, onDismissRequest = { dayExpanded = false }) {
          DayOfWeek.entries.forEach { day ->
            DropdownMenuItem(
              text = { Text(day.toJapaneseLabel()) },
              onClick = {
                selectedDay = day
                dayExpanded = false
              },
            )
          }
        }
      }

      Column {
        TextButton(onClick = { periodExpanded = true }) {
          Text("時限: $selectedPeriod")
        }
        DropdownMenu(expanded = periodExpanded, onDismissRequest = { periodExpanded = false }) {
          (1..6).forEach { period ->
            DropdownMenuItem(
              text = { Text(period.toString()) },
              onClick = {
                selectedPeriod = period
                periodExpanded = false
              },
            )
          }
        }
      }

      Button(
        onClick = {
          val entry = selectedDay to selectedPeriod
          if (!selected.contains(entry)) {
            selected.add(entry)
          }
        },
      ) {
        Text("追加")
      }
    }

    if (selected.isNotEmpty()) {
      Text(
        text = selected.joinToString(", ") { (day, period) -> "${day.toJapaneseLabel()}$period" },
        style = MaterialTheme.typography.bodyMedium,
      )
      TextButton(onClick = { selected.clear() }) {
        Text("クリア")
      }
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun TimetablePickerEmptyPreview() {
  val selected = remember { mutableStateListOf<Pair<DayOfWeek, Int>>() }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      Text("時間割", style = MaterialTheme.typography.titleMedium)
      TimetablePicker(selected)
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun TimetablePickerSelectedPreview() {
  val selected = remember {
    mutableStateListOf(
      DayOfWeek.monday to 1,
      DayOfWeek.wednesday to 3,
    )
  }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      Text("時間割", style = MaterialTheme.typography.titleMedium)
      TimetablePicker(selected)
    }
  }
}

@Composable
private fun LabeledCheckbox(
  label: String,
  checked: Boolean,
  onCheckedChange: (Boolean) -> Unit,
) {
  Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
    Checkbox(checked = checked, onCheckedChange = onCheckedChange)
    Text(label)
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LabeledCheckboxPreview() {
  var checked by remember { mutableStateOf(true) }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LabeledCheckbox(label = "サンプル", checked = checked, onCheckedChange = { checked = it })
    }
  }
}

private fun DayOfWeek.toJapaneseLabel(): String = when (this) {
  DayOfWeek.monday -> "月"
  DayOfWeek.tuesday -> "火"
  DayOfWeek.wednesday -> "水"
  DayOfWeek.thursday -> "木"
  DayOfWeek.friday -> "金"
  DayOfWeek.saturday -> "土"
  DayOfWeek.sunday -> "日"
}

