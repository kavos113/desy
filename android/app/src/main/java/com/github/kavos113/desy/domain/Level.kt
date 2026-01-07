package com.github.kavos113.desy.domain

enum class Level(val value: Int) {
  bachelor1(1),
  bachelor2(2),
  bachelor3(3),
  master1(4),
  master2(5),
  doctor(6);

  companion object {
    fun fromDb(value: Int?): Level? = value?.let { v -> entries.firstOrNull { it.value == v } }
  }
}
