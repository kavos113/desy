package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/kavos113/desy/backend/domain"
)

type timetableRepoStub struct {
	expandFunc func(ctx context.Context) (int, error)
}

func (s *timetableRepoStub) FindByLectureID(int) ([]domain.TimeTable, error) { return nil, nil }
func (s *timetableRepoStub) Create(*domain.TimeTable) error                  { return nil }
func (s *timetableRepoStub) Creates([]domain.TimeTable) error                { return nil }
func (s *timetableRepoStub) Update(*domain.TimeTable) error                  { return nil }
func (s *timetableRepoStub) Delete(int) error                                { return nil }
func (s *timetableRepoStub) ExpandTimetableRanges(ctx context.Context) (int, error) {
	if s.expandFunc != nil {
		return s.expandFunc(ctx)
	}
	return 0, nil
}

func TestTimeTableUsecaseExpandTimetableRanges_Success(t *testing.T) {
	expected := 3
	ctx := context.WithValue(context.Background(), struct{}{}, "marker")

	repo := &timetableRepoStub{
		expandFunc: func(received context.Context) (int, error) {
			if received != ctx {
				t.Fatalf("context mismatch")
			}
			return expected, nil
		},
	}

	uc := NewTimeTableUsecase(repo)
	got, err := uc.ExpandTimetableRanges(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Fatalf("unexpected result: got %d want %d", got, expected)
	}
}

func TestTimeTableUsecaseExpandTimetableRanges_Error(t *testing.T) {
	expectedErr := errors.New("boom")
	uc := NewTimeTableUsecase(&timetableRepoStub{
		expandFunc: func(context.Context) (int, error) {
			return 0, expectedErr
		},
	})

	_, err := uc.ExpandTimetableRanges(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestTimeTableUsecaseExpandTimetableRanges_RepositoryNotInitialized(t *testing.T) {
	uc := NewTimeTableUsecase(nil)
	if _, err := uc.ExpandTimetableRanges(context.Background()); err == nil {
		t.Fatal("expected error but got nil")
	}

	var nilUsecase *timetableUsecase
	if _, err := nilUsecase.ExpandTimetableRanges(context.Background()); err == nil {
		t.Fatal("expected error when usecase is nil")
	}
}
