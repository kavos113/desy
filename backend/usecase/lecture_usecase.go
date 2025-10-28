package usecase

import (
	"context"
	"errors"

	"github.com/kavos113/desy/backend/domain"
)

// LectureUsecase defines application logic related to lectures.
type LectureUsecase interface {
	SearchLectures(query domain.SearchQuery) ([]domain.LectureSummary, error)
	GetLectureDetails(lectureID int) (*domain.Lecture, error)
	MigrateRelatedCourses(ctx context.Context) (int, error)
}

// lectureUsecase is a concrete implementation of LectureUsecase.
type lectureUsecase struct {
	lectureRepo domain.LectureRepository
}

// NewLectureUsecase creates a new lecture usecase instance.
func NewLectureUsecase(lectureRepo domain.LectureRepository) LectureUsecase {
	return &lectureUsecase{
		lectureRepo: lectureRepo,
	}
}

// SearchLectures retrieves lecture summaries based on the provided query.
func (uc *lectureUsecase) SearchLectures(query domain.SearchQuery) ([]domain.LectureSummary, error) {
	if uc.lectureRepo == nil {
		return nil, errors.New("lecture repository is not initialized")
	}

	return uc.lectureRepo.Search(query)
}

// GetLectureDetails retrieves a full lecture aggregate by its identifier.
func (uc *lectureUsecase) GetLectureDetails(lectureID int) (*domain.Lecture, error) {
	if uc.lectureRepo == nil {
		return nil, errors.New("lecture repository is not initialized")
	}

	return uc.lectureRepo.FindByID(lectureID)
}

// MigrateRelatedCourses resolves related course IDs based on stored course codes.
func (uc *lectureUsecase) MigrateRelatedCourses(ctx context.Context) (int, error) {
	if uc.lectureRepo == nil {
		return 0, errors.New("lecture repository is not initialized")
	}

	return uc.lectureRepo.MigrateRelatedCourses(ctx)
}
