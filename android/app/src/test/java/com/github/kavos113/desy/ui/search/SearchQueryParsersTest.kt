package com.github.kavos113.desy.ui.search

import org.junit.Assert.assertEquals
import org.junit.Test

class SearchQueryParsersTest {
  @Test
  fun parseKeywordInput_splitsBySpacesAndComma() {
    assertEquals(listOf("a", "b", "c"), parseKeywordInput("a b,c"))
  }

  @Test
  fun parseDepartmentInput_trimsAndFiltersEmpty() {
    assertEquals(listOf("理学院", "工学院"), parseDepartmentInput(" 理学院  工学院 "))
  }

  @Test
  fun parseYearInput_returnsZeroOnInvalid() {
    assertEquals(0, parseYearInput(""))
    assertEquals(0, parseYearInput("abc"))
    assertEquals(2025, parseYearInput("2025"))
  }
}
