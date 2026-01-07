package com.github.kavos113.desy.repository.db

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import com.github.kavos113.desy.repository.db.dao.LectureDao
import com.github.kavos113.desy.repository.db.entity.LectureEntity
import com.github.kavos113.desy.repository.db.entity.LectureKeywordEntity
import com.github.kavos113.desy.repository.db.entity.LecturePlanEntity
import com.github.kavos113.desy.repository.db.entity.LectureTeacherCrossRef
import com.github.kavos113.desy.repository.db.entity.RelatedCourseCodeEntity
import com.github.kavos113.desy.repository.db.entity.RelatedCourseEntity
import com.github.kavos113.desy.repository.db.entity.RoomEntity
import com.github.kavos113.desy.repository.db.entity.TeacherEntity
import com.github.kavos113.desy.repository.db.entity.TimetableEntity
import java.io.File

@Database(
  entities = [
    LectureEntity::class,
    TeacherEntity::class,
    LectureTeacherCrossRef::class,
    RoomEntity::class,
    TimetableEntity::class,
    LecturePlanEntity::class,
    LectureKeywordEntity::class,
    RelatedCourseEntity::class,
    RelatedCourseCodeEntity::class,
  ],
  version = 1,
  exportSchema = false,
)
abstract class DesyDatabase : RoomDatabase() {
  abstract fun lectureDao(): LectureDao

  companion object {
    /**
     * デスクトップ版で生成したSQLite DBを配置しておけば、それを読み込みます。
     * dbFileを指定しない場合はアプリ内の `desy.db` を開きます。
     */
    fun open(context: Context, dbFile: File? = null): DesyDatabase {
      val builder = Room.databaseBuilder(context, DesyDatabase::class.java, "desy.db")
      if (dbFile != null) {
        builder.createFromFile(dbFile)
      }
      return builder.build()
    }
  }
}
