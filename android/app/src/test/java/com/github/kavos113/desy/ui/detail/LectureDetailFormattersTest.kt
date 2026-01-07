package com.github.kavos113.desy.ui.detail

import org.junit.Assert.assertEquals
import org.junit.Test

class LectureDetailFormattersTest {
  @Test
  fun splitIntoLines_handlesNullAndBlankLines() {
    assertEquals(emptyList<String>(), splitIntoLines(null))

    val input = "line1\n\nline2\r\n   \r\nline3  "
    assertEquals(listOf("line1", "line2", "line3"), splitIntoLines(input))
  }
}
