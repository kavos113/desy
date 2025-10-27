package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/scraper"
)

// Fetcher describes the ability to load HTML content from a URL.
type Fetcher interface {
	Fetch(ctx context.Context, url string) (io.ReadCloser, error)
}

// ScraperUsecase orchestrates scraping workflow and persistence.
type ScraperUsecase interface {
	ScrapeCourseList(ctx context.Context, listURL, baseURL string) ([]scraper.CourseListItem, error)
	ScrapeCourseListAndSave(ctx context.Context, listURL, baseURL string) ([]domain.Lecture, error)
	ScrapeCourseDetail(ctx context.Context, detailURL string) (*domain.Lecture, error)
	ScrapeCourseDetailAndSave(ctx context.Context, detailURL string) (*domain.Lecture, error)
}

type scraperUsecase struct {
	fetcher     Fetcher
	lectureRepo domain.LectureRepository
}

// NewScraperUsecase constructs a scraper usecase instance.
func NewScraperUsecase(fetcher Fetcher, lectureRepo domain.LectureRepository) ScraperUsecase {
	return &scraperUsecase{
		fetcher:     fetcher,
		lectureRepo: lectureRepo,
	}
}

// ScrapeCourseList retrieves course list entries from the specified URL.
func (uc *scraperUsecase) ScrapeCourseList(ctx context.Context, listURL, baseURL string) ([]scraper.CourseListItem, error) {
	if uc.fetcher == nil {
		return nil, errors.New("scraper fetcher is not initialized")
	}

	if strings.TrimSpace(listURL) == "" {
		return nil, errors.New("list url is required")
	}

	reader, err := uc.fetcher.Fetch(ctx, listURL)
	if err != nil {
		return nil, fmt.Errorf("fetch list %s: %w", listURL, err)
	}
	defer reader.Close()

	items, err := scraper.ParseCourseList(reader, baseURL)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// ScrapeCourseListAndSave scrapes a course list, expands each detail, and persists the result.
func (uc *scraperUsecase) ScrapeCourseListAndSave(ctx context.Context, listURL, baseURL string) ([]domain.Lecture, error) {
	if uc.lectureRepo == nil {
		return nil, errors.New("lecture repository is not initialized")
	}

	items, err := uc.ScrapeCourseList(ctx, listURL, baseURL)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []domain.Lecture{}, nil
	}

	lectures := make([]domain.Lecture, 0, len(items))
	seen := make(map[string]struct{})

	for _, item := range items {
		if ctx != nil && ctx.Err() != nil {
			return nil, ctx.Err()
		}

		detailURL := strings.TrimSpace(item.DetailURL)
		if detailURL == "" {
			continue
		}
		if _, ok := seen[detailURL]; ok {
			continue
		}
		seen[detailURL] = struct{}{}

		lecture, err := uc.ScrapeCourseDetail(ctx, detailURL)
		if err != nil {
			return nil, err
		}
		if lecture == nil {
			continue
		}
		if lecture.Code == "" {
			lecture.Code = item.Code
		}
		if lecture.Title == "" {
			lecture.Title = item.Title
		}
		lectures = append(lectures, *lecture)
	}

	if len(lectures) == 0 {
		return []domain.Lecture{}, nil
	}

	if err := uc.lectureRepo.Creates(lectures); err != nil {
		return nil, err
	}

	return lectures, nil
}

// ScrapeCourseDetail retrieves a single lecture aggregate from a detail page.
func (uc *scraperUsecase) ScrapeCourseDetail(ctx context.Context, detailURL string) (*domain.Lecture, error) {
	if uc.fetcher == nil {
		return nil, errors.New("scraper fetcher is not initialized")
	}

	if strings.TrimSpace(detailURL) == "" {
		return nil, errors.New("detail url is required")
	}

	reader, err := uc.fetcher.Fetch(ctx, detailURL)
	if err != nil {
		return nil, fmt.Errorf("fetch detail %s: %w", detailURL, err)
	}
	defer reader.Close()

	lecture, err := scraper.ParseCourseDetail(reader, detailURL)
	if err != nil {
		return nil, err
	}

	return lecture, nil
}

// ScrapeCourseDetailAndSave scrapes a detail page and persists the aggregate.
func (uc *scraperUsecase) ScrapeCourseDetailAndSave(ctx context.Context, detailURL string) (*domain.Lecture, error) {
	if uc.lectureRepo == nil {
		return nil, errors.New("lecture repository is not initialized")
	}

	lecture, err := uc.ScrapeCourseDetail(ctx, detailURL)
	if err != nil {
		return nil, err
	}
	if lecture == nil {
		return nil, errors.New("scraped lecture is nil")
	}

	if err := uc.lectureRepo.Create(lecture); err != nil {
		return nil, err
	}

	return lecture, nil
}
