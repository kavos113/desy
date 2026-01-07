package com.github.kavos113.desy.repository.db.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "rooms")
data class RoomEntity(
  @PrimaryKey
  @ColumnInfo(name = "id")
  val id: Int,

  @ColumnInfo(name = "name")
  val name: String,
)
