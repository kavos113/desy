package com.github.kavos113.desy.domain

enum class DayOfWeek {
  monday,
  tuesday,
  wednesday,
  thursday,
  friday,
  saturday,
  sunday;

  companion object {
    fun fromDb(value: String?): DayOfWeek? = value
      ?.trim()
      ?.takeIf { it.isNotEmpty() }
      ?.let { token -> entries.firstOrNull { it.name == token } }
  }
}
