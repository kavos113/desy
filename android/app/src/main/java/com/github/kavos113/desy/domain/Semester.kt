package com.github.kavos113.desy.domain

enum class Semester {
  spring,
  summer,
  fall,
  winter;

  companion object {
    fun fromDb(value: String?): Semester? = value
      ?.trim()
      ?.takeIf { it.isNotEmpty() }
      ?.let { token -> entries.firstOrNull { it.name == token } }
  }
}
