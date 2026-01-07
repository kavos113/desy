package com.github.kavos113.desy.ui.search

private val KEYWORD_SEPARATOR = Regex("[\\s,\\u3001\\u3002、;；]+")

internal fun parseKeywordInput(value: String): List<String> {
  return value
    .split(KEYWORD_SEPARATOR)
    .map { it.trim() }
    .filter { it.isNotEmpty() }
}

internal fun parseDepartmentInput(value: String): List<String> {
  return value
    .split(KEYWORD_SEPARATOR)
    .map { it.trim() }
    .filter { it.isNotEmpty() }
}

internal fun parseYearInput(value: String): Int {
  val trimmed = value.trim()
  if (trimmed.isEmpty()) return 0
  return trimmed.toIntOrNull() ?: 0
}
