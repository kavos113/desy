package com.github.kavos113.desy.di

import android.content.Context
import com.github.kavos113.desy.domain.LectureRepository
import com.github.kavos113.desy.repository.RoomLectureRepository
import com.github.kavos113.desy.repository.db.DesyDatabase
import dagger.Binds
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class AppModule {
  companion object {
    @Provides
    @Singleton
    fun provideDatabase(
      @ApplicationContext context: Context,
    ): DesyDatabase {
      // ここでは内部DB(desy.db)を開く。
      // デスクトップ版のdbを読みたい場合は createFromFile を使う open の引数を差し替える。
      return DesyDatabase.open(context, dbFile = null)
    }
  }

  @Binds
  @Singleton
  abstract fun bindLectureRepository(impl: RoomLectureRepository): LectureRepository
}
