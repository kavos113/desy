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
      return DesyDatabase.open(context)
    }
  }

  @Binds
  @Singleton
  abstract fun bindLectureRepository(impl: RoomLectureRepository): LectureRepository
}
