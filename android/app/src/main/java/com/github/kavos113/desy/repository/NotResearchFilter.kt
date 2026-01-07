package com.github.kavos113.desy.repository

import com.github.kavos113.desy.domain.LectureSummary

/** backendの filterNotResearch/checkIsResearchLecture を移植 */
object NotResearchFilter {
  private val researchWords = listOf(
    "研究",
    "卒業",
    "修士",
    "博士",
    "実験",
    "実習",
    "輪講",
    "ゼミ",
    "演習",
    "課題研究",
    "特別研究",
    "卒業研究",
    "卒業制作",
    "卒業設計",
    "卒業論文",
    "修士研究",
    "修士論文",
    "博士研究",
    "博士論文",
    "博士論文指導",
    "博士論文研究",
    "博士論文演習",
    "博士論文特別研究",
    "博士論文指導演習",
    "博士論文指導特別研究",
    "博士論文指導特別演習",
    "博士論文指導特別研究演習",
    "博士論文指導演習特別研究",
    "博士論文研究指導",
    "博士論文研究指導演習",
    "博士論文研究指導特別研究",
    "博士論文研究指導特別演習",
    "博士論文研究指導特別研究演習",
    "博士論文研究指導演習特別研究",
    "博士論文研究指導演習特別",
    "博士論文研究指導特別研究特別",
    "博士論文研究指導特別演習特別",
    "博士論文研究指導特別研究演習特別",
    "博士論文研究指導演習特別研究特別",
    "講究",
    "B2D",
    "リカレント研修",
    "オフキャンパスプロジェクト",
    "国際派遣プロジェクト",
    "インターンシップ",
    "学外研修",
    "論文研究計画論",
    "物理学プレゼンテーション実践",
    "数理・計算科学プレゼンテーション実践",
    "国際プレゼンテーション",
    "エンジニアリングデザインプレゼンテーション実践",
    "チュートリアル",
    "留学",
    "国際研究",
    "キャリアディベロップメント",
    "キャリア開発",
    "キャリア特別",
    "派遣プロジェクト",
    "企画実践",
    "先端研究",
  )

  private val researchTeacherWords = listOf(
    "教員",
  )

  fun filter(lectures: List<LectureSummary>): List<LectureSummary> {
    return lectures.filterNot { isResearchLecture(it.title, it.teachers.map { t -> t.name }) }
  }

  private fun isResearchLecture(title: String, teacherNames: List<String>): Boolean {
    if (researchWords.any { title.contains(it) }) return true
    return teacherNames.any { name -> researchTeacherWords.any { w -> name.contains(w) } }
  }
}
