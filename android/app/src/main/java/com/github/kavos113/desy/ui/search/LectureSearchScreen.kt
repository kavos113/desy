package com.github.kavos113.desy.ui.search

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.clickable
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.Checkbox
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuAnchorType
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.RectangleShape
import androidx.compose.ui.Alignment
import androidx.compose.ui.unit.Dp
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
  var keywordsText by remember { mutableStateOf("") }
  var filterNotResearch by remember { mutableStateOf(false) }

  val selectedLevels = remember { mutableStateListOf<Level>() }
  val selectedSemesters = remember { mutableStateListOf<Semester>() }
  val selectedTimetables = remember { mutableStateListOf<Pair<DayOfWeek, Int>>() }
  val selectedDepartments = remember { mutableStateListOf<String>() }

  LazyColumn(
    modifier = modifier
      .fillMaxSize()
      .padding(12.dp),
    verticalArrangement = Arrangement.spacedBy(12.dp),
  ) {
    item {
      Row(modifier = Modifier.fillMaxWidth()) {
        Text("検索", style = MaterialTheme.typography.titleLarge)
        Spacer(modifier = Modifier.weight(1f))
        TextButton(onClick = onCancel) {
          Text("戻る")
        }
      }
    }

    item {
      OutlinedTextField(
        value = title,
        onValueChange = { title = it },
        modifier = Modifier.fillMaxWidth(),
        label = { Text("講義名") },
        singleLine = true,
      )
    }

    item {
      OutlinedTextField(
        value = teacherName,
        onValueChange = { teacherName = it },
        modifier = Modifier.fillMaxWidth(),
        label = { Text("教員名") },
        singleLine = true,
      )
    }

    item {
      OutlinedTextField(
        value = room,
        onValueChange = { room = it },
        modifier = Modifier.fillMaxWidth(),
        label = { Text("教室") },
        singleLine = true,
      )
    }

    item {
      OutlinedTextField(
        value = yearText,
        onValueChange = { yearText = it },
        modifier = Modifier.fillMaxWidth(),
        label = { Text("年度") },
        singleLine = true,
        placeholder = { Text("例: 2025") },
      )
    }

    item {
      DepartmentMultiSelect(
        label = "開講元",
        options = DepartmentOptions.mobileDepartments,
        selected = selectedDepartments,
      )
    }

    item {
      OutlinedTextField(
        value = keywordsText,
        onValueChange = { keywordsText = it },
        modifier = Modifier.fillMaxWidth(),
        label = { Text("キーワード（複数は空白/カンマ区切り）") },
      )
    }

    item {
      HorizontalDivider()
    }

    item {
      Text("学年", style = MaterialTheme.typography.titleMedium)
    }
    item {
      LevelCheckboxRow(selectedLevels)
    }

    item {
      Spacer(Modifier.height(4.dp))
    }

    item {
      Text("学期(Quarter)", style = MaterialTheme.typography.titleMedium)
    }
    item {
      SemesterCheckboxRow(selectedSemesters)
    }

    item {
      HorizontalDivider()
    }

    item {
      Text("時間割", style = MaterialTheme.typography.titleMedium)
    }
    item {
      TimetablePicker(
        selected = selectedTimetables,
      )
    }

    item {
      HorizontalDivider()
    }

    item {
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
              departments = selectedDepartments.toList(),
              keywords = parseKeywordInput(keywordsText),
              semesters = selectedSemesters.toList(),
              timetables = selectedTimetables.map { (day, period) ->
                listOf(
                  TimeTable(lectureId = 0, dayOfWeek = day, period = 2 * period - 1),
                  TimeTable(lectureId = 0, dayOfWeek = day, period = 2 * period)
                )
              }.flatten(),
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
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun DepartmentMultiSelect(
  label: String,
  options: List<String>,
  selected: MutableList<String>,
  modifier: Modifier = Modifier,
) {
  var expanded by remember { mutableStateOf(false) }

  ExposedDropdownMenuBox(
    expanded = expanded,
    onExpandedChange = { expanded = !expanded },
    modifier = modifier.fillMaxWidth()
  ) {
    OutlinedTextField(
      value = selected.joinToString(", "),
      onValueChange = {},
      modifier = Modifier
        .fillMaxWidth()
        .menuAnchor(ExposedDropdownMenuAnchorType.PrimaryNotEditable),
      readOnly = true,
      label = { Text(label) },
      placeholder = { Text("選択してください") },
      trailingIcon = {
        ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded)
      },
      colors = ExposedDropdownMenuDefaults.outlinedTextFieldColors()
    )

    ExposedDropdownMenu(
      expanded = expanded,
      onDismissRequest = { expanded = false },
      modifier = Modifier.heightIn(max = 500.dp)
    ) {
      options.forEach { option ->
        val checked = selected.contains(option)
        DropdownMenuItem(
          text = {
            Row(
              modifier = Modifier.fillMaxWidth(),
              horizontalArrangement = Arrangement.spacedBy(8.dp),
              verticalAlignment = Alignment.CenterVertically,
            ) {
              Checkbox(
                checked = checked,
                onCheckedChange = null,
              )
              Text(option)
            }
          },
          onClick = {
            if (checked) {
              selected.remove(option)
            } else {
              selected.add(option)
            }
          },
        )
      }
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun DepartmentMultiSelectPreview() {
  val selected = remember { mutableStateListOf("理学院", "情報工学系") }
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
      DepartmentMultiSelect(
        label = "開講元",
        options = DepartmentOptions.mobileDepartments,
        selected = selected,
      )
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
  val days = listOf(
    DayOfWeek.monday,
    DayOfWeek.tuesday,
    DayOfWeek.wednesday,
    DayOfWeek.thursday,
    DayOfWeek.friday,
  )
  val periods = (1..5).toList()

  fun isChecked(day: DayOfWeek, period: Int): Boolean {
    return selected.contains(day to period)
  }

  fun toggle(day: DayOfWeek, period: Int) {
    val entry = day to period
    if (selected.contains(entry)) {
      selected.remove(entry)
    } else {
      selected.add(entry)
    }
  }

  fun toggleDay(day: DayOfWeek) {
    periods.forEach { period ->
      toggle(day, period)
    }
  }

  fun togglePeriod(period: Int) {
    days.forEach { day ->
      toggle(day, period)
    }
  }

  val cellSize: Dp = 28.dp
  val headerWidth: Dp = 48.dp

  Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
    Column {
      // Header row
      Row {
        Box(
          modifier = Modifier
            .size(width = headerWidth, height = cellSize),
        )
        days.forEach { day ->
          Box(
            modifier = Modifier
              .size(cellSize)
              .border(1.dp, MaterialTheme.colorScheme.outline)
              .clip(RectangleShape)
              .clickable { toggleDay(day) },
            contentAlignment = Alignment.Center,
          ) {
            Text(day.toJapaneseLabel(), style = MaterialTheme.typography.labelLarge)
          }
        }
      }

      // Body
      periods.forEach { period ->
        Row {
          Box(
            modifier = Modifier
              .size(width = headerWidth, height = cellSize)
              .border(1.dp, MaterialTheme.colorScheme.outline)
              .clip(RectangleShape)
              .clickable { togglePeriod(period) },
            contentAlignment = Alignment.Center,
          ) {
            Text(period.toString(), style = MaterialTheme.typography.labelLarge)
          }

          days.forEach { day ->
            val checked = isChecked(day, period)
            Box(
              modifier = Modifier
                .size(cellSize)
                .border(1.dp, MaterialTheme.colorScheme.outline)
                .background(
                  if (checked) MaterialTheme.colorScheme.primary.copy(alpha = 0.25f)
                  else MaterialTheme.colorScheme.surface
                )
                .clip(RectangleShape)
                .clickable { toggle(day, period) },
            )
          }
        }
      }
    }

    if (selected.isNotEmpty()) {
      val dayOrder = mapOf(
        DayOfWeek.monday to 1,
        DayOfWeek.tuesday to 2,
        DayOfWeek.wednesday to 3,
        DayOfWeek.thursday to 4,
        DayOfWeek.friday to 5,
        DayOfWeek.saturday to 6,
        DayOfWeek.sunday to 7,
      )
      Text(
        text = selected
          .sortedWith(compareBy({ dayOrder[it.first] ?: 99 }, { it.second }))
          .joinToString(", ") { (day, period) -> "${day.toJapaneseLabel()}$period" },
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
  Row(
    horizontalArrangement = Arrangement.spacedBy(6.dp),
    verticalAlignment = Alignment.CenterVertically,
    modifier = Modifier
      .height(24.dp)
  ) {
    Checkbox(
      checked = checked,
      onCheckedChange = onCheckedChange,
    )
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

