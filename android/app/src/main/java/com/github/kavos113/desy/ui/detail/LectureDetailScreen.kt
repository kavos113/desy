package com.github.kavos113.desy.ui.detail

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.github.kavos113.desy.domain.DayOfWeek
import com.github.kavos113.desy.domain.Lecture
import com.github.kavos113.desy.domain.LecturePlan
import com.github.kavos113.desy.domain.Semester
import com.github.kavos113.desy.domain.Teacher
import com.github.kavos113.desy.domain.TimeTable
import com.github.kavos113.desy.ui.theme.DesyTheme
import com.github.kavos113.desy.ui.viewmodel.LectureDetailUiState

@Composable
fun LectureDetailScreen(
  uiState: LectureDetailUiState,
  onBack: () -> Unit,
  onSelectRelatedLecture: (Int) -> Unit,
  modifier: Modifier = Modifier,
) {
  val listState = rememberLazyListState()
  LaunchedEffect(uiState.lectureId) {
    if (uiState.lectureId != null) {
      listState.scrollToItem(0)
    }
  }

  Column(
    modifier = modifier
      .fillMaxSize()
      .padding(12.dp),
    verticalArrangement = Arrangement.spacedBy(12.dp),
  ) {
    Row(modifier = Modifier.fillMaxWidth()) {
      TextButton(onClick = onBack) {
        Text("戻る")
      }
      Spacer(modifier = Modifier.weight(1f))
      Text("講義詳細", style = MaterialTheme.typography.titleLarge)
    }

    when {
      uiState.isLoading -> {
        Text("読み込み中…")
      }

      uiState.errorMessage != null -> {
        Text(uiState.errorMessage)
      }

      uiState.lecture == null -> {
        Text("講義の詳細が選択されていません。")
      }

      else -> {
        LazyColumn(
          modifier = Modifier.fillMaxWidth(),
          state = listState,
          verticalArrangement = Arrangement.spacedBy(12.dp),
        ) {
          item {
            LectureTitleSection(uiState.lecture)
          }

          item {
            LectureMetaSection(lecture = uiState.lecture)
          }

          item { HorizontalDivider() }

          item {
            LectureTextSection(title = "講義の概要とねらい", content = uiState.lecture.abstractText)
          }

          item {
            LecturePlansSection(plans = uiState.lecture.lecturePlans)
          }

          item {
            LectureTextSection(title = "到達目標", content = uiState.lecture.goal)
          }

          if (!uiState.lecture.experience.isNullOrBlank()) {
            item {
              LectureTextSection(title = "実務経験のある教員による授業", content = "あり")
            }
          }

          item {
            LectureKeywordsSection(uiState.lecture.keywords)
          }

          item {
            LectureListSection(title = "教科書", content = uiState.lecture.textbook)
          }

          item {
            LectureListSection(title = "参考書・講義資料等", content = uiState.lecture.referenceBook)
          }

          if (!uiState.lecture.flow.isNullOrBlank()) {
            item { LectureTextSection(title = "授業の進め方", content = uiState.lecture.flow) }
          }
          if (!uiState.lecture.outOfClassWork.isNullOrBlank()) {
            item { LectureTextSection(title = "授業時間外学修（予習・復習等）", content = uiState.lecture.outOfClassWork) }
          }
          if (!uiState.lecture.assessment.isNullOrBlank()) {
            item { LectureTextSection(title = "成績評価の基準及び方法", content = uiState.lecture.assessment) }
          }

          item {
            LectureRelatedSection(
              relatedCourses = uiState.relatedCourses,
              onSelectRelatedLecture = onSelectRelatedLecture,
            )
          }

          if (!uiState.lecture.prerequisite.isNullOrBlank()) {
            item { LectureTextSection(title = "履修の条件", content = uiState.lecture.prerequisite) }
          }
          if (!uiState.lecture.note.isNullOrBlank()) {
            item { LectureTextSection(title = "その他", content = uiState.lecture.note) }
          }
          if (!uiState.lecture.contact.isNullOrBlank()) {
            item { LectureTextSection(title = "連絡先", content = uiState.lecture.contact) }
          }
          if (!uiState.lecture.officeHours.isNullOrBlank()) {
            item { LectureTextSection(title = "オフィスアワー", content = uiState.lecture.officeHours) }
          }
        }
      }
    }
  }
}

@Composable
private fun LectureTitleSection(lecture: Lecture) {
  Column(verticalArrangement = Arrangement.spacedBy(4.dp)) {
    Text(lecture.title, style = MaterialTheme.typography.titleLarge, fontWeight = FontWeight.SemiBold)
    Text(lecture.englishTitle.orEmpty(), style = MaterialTheme.typography.bodyMedium)
  }
}

@Composable
private fun LectureMetaSection(lecture: Lecture) {
  val timetableText = formatTimetablesWithRoom(lecture.timetables)
  val semesterText = formatSemesters(lecture.timetables)
  val teacherText = formatTeachers(lecture.teachers)

  Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
    MetaRow(label = "開講元", value = lecture.department.orEmpty())
    MetaRow(label = "担当教員", value = teacherText)
    MetaRow(label = "授業形態", value = lecture.lectureType?.name.orEmpty())
    MetaRow(label = "曜日・時限(講義室)", value = timetableText)

    val yearPart = lecture.year?.let { "${it}年" }
    val openPeriod = listOfNotNull(yearPart, semesterText.takeIf { it.isNotBlank() }).joinToString(" ")
    MetaRow(label = "開講時期", value = openPeriod)

    Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(16.dp)) {
      Column(modifier = Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(2.dp)) {
        Text("科目コード", style = MaterialTheme.typography.labelLarge)
        Text(lecture.code.orEmpty(), style = MaterialTheme.typography.bodyMedium)
      }
      Column(modifier = Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(2.dp)) {
        Text("単位数", style = MaterialTheme.typography.labelLarge)
        Text(lecture.credit?.toString().orEmpty(), style = MaterialTheme.typography.bodyMedium)
      }
    }

    Row(modifier = Modifier.fillMaxWidth()) {
      Column(modifier = Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(2.dp)) {
        Text("言語", style = MaterialTheme.typography.labelLarge)
        Text(lecture.language.orEmpty(), style = MaterialTheme.typography.bodyMedium)
      }
    }
  }
}

@Composable
private fun MetaRow(label: String, value: String) {
  Column(verticalArrangement = Arrangement.spacedBy(2.dp)) {
    Text(label, style = MaterialTheme.typography.labelLarge)
    Text(if (value.isBlank()) "" else value, style = MaterialTheme.typography.bodyMedium)
  }
}

@Composable
private fun LectureTextSection(title: String, content: String?) {
  Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
    Text(title, style = MaterialTheme.typography.titleMedium)
    val lines = splitIntoLines(content)
    if (lines.isEmpty()) {
      Text("", style = MaterialTheme.typography.bodyMedium)
    } else {
      lines.forEach { line ->
        Text(line, style = MaterialTheme.typography.bodyMedium)
      }
    }
  }
}

@Composable
private fun LectureKeywordsSection(keywords: List<String>) {
  Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
    Text("キーワード", style = MaterialTheme.typography.titleMedium)
    Text(keywords.joinToString(", "), style = MaterialTheme.typography.bodyMedium)
  }
}

@Composable
private fun LectureListSection(title: String, content: String?) {
  Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
    Text(title, style = MaterialTheme.typography.titleMedium)
    val items = splitIntoLines(content)
    if (items.isEmpty()) {
      Text("", style = MaterialTheme.typography.bodyMedium)
    } else {
      items.forEach { item ->
        Text("・$item", style = MaterialTheme.typography.bodyMedium)
      }
    }
  }
}

@Composable
private fun LecturePlansSection(plans: List<LecturePlan>) {
  Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
    Text("授業計画・課題", style = MaterialTheme.typography.titleMedium)

    if (plans.isEmpty()) {
      Text("", style = MaterialTheme.typography.bodyMedium)
      return
    }

    Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
      Text("回", modifier = Modifier.weight(0.12f), style = MaterialTheme.typography.labelLarge)
      Text("授業計画", modifier = Modifier.weight(0.44f), style = MaterialTheme.typography.labelLarge)
      Text("課題", modifier = Modifier.weight(0.44f), style = MaterialTheme.typography.labelLarge)
    }

    plans.sortedBy { it.count }.forEach { plan ->
      Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
        Text("第${plan.count}回", modifier = Modifier.weight(0.12f), style = MaterialTheme.typography.bodySmall)
        Text(plan.plan.orEmpty(), modifier = Modifier.weight(0.44f), style = MaterialTheme.typography.bodySmall)
        Text(plan.assignment.orEmpty(), modifier = Modifier.weight(0.44f), style = MaterialTheme.typography.bodySmall)
      }
    }
  }
}

@Composable
private fun LectureRelatedSection(
  relatedCourses: List<RelatedCourseEntry>,
  onSelectRelatedLecture: (Int) -> Unit,
) {
  Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
    Text("関連する科目", style = MaterialTheme.typography.titleMedium)

    if (relatedCourses.isEmpty()) {
      Text("関連科目の情報がありません。", style = MaterialTheme.typography.bodyMedium)
      return
    }

    Column(verticalArrangement = Arrangement.spacedBy(4.dp)) {
      relatedCourses.forEach { course ->
        val clickable = course.id != null
        val label = if (!course.title.isNullOrBlank()) {
          "${course.title} (${course.code})"
        } else {
          course.code
        }

        if (clickable) {
          TextButton(onClick = { onSelectRelatedLecture(course.id!!) }) {
            Text(label)
          }
        } else {
          val text = if (!course.title.isNullOrBlank()) {
            "${course.code} / ${course.title}"
          } else {
            course.code
          }
          Text(text, style = MaterialTheme.typography.bodyMedium)
        }
      }
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureDetailScreenPreview() {
  val lecture = sampleLecture()
  DesyTheme {
    LectureDetailScreen(
      uiState = LectureDetailUiState(
        lectureId = lecture.id,
        lecture = lecture,
        relatedCourses = listOf(
          RelatedCourseEntry(code = "CS102"),
          RelatedCourseEntry(code = "CS201", id = 2, title = "データ構造"),
        ),
      ),
      onBack = {},
      onSelectRelatedLecture = {},
    )
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureTitleSectionPreview() {
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LectureTitleSection(sampleLecture())
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureMetaSectionPreview() {
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
      LectureMetaSection(sampleLecture())
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LecturePlansSectionPreview() {
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LecturePlansSection(sampleLecture().lecturePlans)
    }
  }
}

@Preview(showBackground = true, widthDp = 420)
@Composable
private fun LectureRelatedSectionPreview() {
  DesyTheme {
    Column(modifier = Modifier.padding(12.dp)) {
      LectureRelatedSection(
        relatedCourses = listOf(
          RelatedCourseEntry(code = "CS101"),
          RelatedCourseEntry(code = "CS201", id = 10, title = "アルゴリズム"),
        ),
        onSelectRelatedLecture = {},
      )
    }
  }
}

private fun sampleLecture(): Lecture {
  return Lecture(
    id = 1,
    university = "Sample University",
    title = "プログラミング基礎",
    englishTitle = "Introduction to Programming",
    department = "情報学部",
    code = "CS101",
    year = 2025,
    openTerm = "1Q",
    credit = 2,
    language = "日本語",
    abstractText = "本講義ではプログラミングの基礎を学ぶ。\n変数・制御構文・関数を扱う。",
    goal = "簡単なプログラムを自力で書けるようになる。",
    textbook = "教科書A\n教科書B",
    referenceBook = "参考書X\n参考書Y",
    assessment = "課題・試験で評価する。",
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
    teachers = listOf(
      Teacher(id = 1, name = "山田 太郎", url = null),
      Teacher(id = 2, name = "佐藤 花子", url = null),
    ),
    lecturePlans = listOf(
      LecturePlan(count = 1, plan = "イントロ", assignment = "環境構築"),
      LecturePlan(count = 2, plan = "変数と式", assignment = "演習1"),
    ),
    keywords = listOf("プログラミング", "基礎"),
    relatedCourseCodes = listOf("CS102"),
    relatedCourses = listOf(2, 3),
  )
}
