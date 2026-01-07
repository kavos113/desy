package com.github.kavos113.desy.domain

enum class LectureType {
  offline,
  live,
  hyflex,
  ondemand,
  other;

  companion object {
    fun fromDb(value: String?): LectureType? =
      value?.trim()?.takeIf { it.isNotEmpty() }?.let { token ->
        entries.firstOrNull { it.name == token }
      }
  }
}
