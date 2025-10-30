package usecase

import (
	"context"
	"errors"

	"github.com/kavos113/desy/backend/domain"
)

type TimeTableUsecase interface {
	ExpandTimetableRanges(ctx context.Context) (int, error)
}

type timetableUsecase struct {
	timetableRepo domain.TimeTableRepository
}

func NewTimeTableUsecase(timetableRepo domain.TimeTableRepository) TimeTableUsecase {
	return &timetableUsecase{
		timetableRepo: timetableRepo,
	}
}

func (uc *timetableUsecase) ExpandTimetableRanges(ctx context.Context) (int, error) {
	if uc == nil || uc.timetableRepo == nil {
		return 0, errors.New("timetable repository is not initialized")
	}

	return uc.timetableRepo.ExpandTimetableRanges(ctx)
}
