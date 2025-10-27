package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/scraper"
)

const defaultScrapeDelay = 2 * time.Second

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
	ScrapeTopPageAndSave(ctx context.Context, year int) ([]domain.Lecture, error)
}

type scraperUsecase struct {
	fetcher     Fetcher
	lectureRepo domain.LectureRepository
	parser      scraper.Parser
	delay       time.Duration
}

// NewScraperUsecase constructs a scraper usecase instance.
func NewScraperUsecase(fetcher Fetcher, lectureRepo domain.LectureRepository, parser scraper.Parser, delay time.Duration) ScraperUsecase {
	if parser == nil {
		parser = scraper.NewParser()
	}
	if delay < 0 {
		delay = defaultScrapeDelay
	}
	return &scraperUsecase{
		fetcher:     fetcher,
		lectureRepo: lectureRepo,
		parser:      parser,
		delay:       delay,
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

	items, err := uc.parser.ParseCourseList(reader, baseURL)
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
	firstFetch := true

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

		log.Printf("Scraping detail page: %s %s", item.Code, item.Title)

		if !firstFetch {
			if err := uc.sleep(ctx); err != nil {
				return nil, err
			}
		}
		firstFetch = false

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

	lecture, err := uc.parser.ParseCourseDetail(reader, detailURL)
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

// ScrapeTopPageAndSave fetches the top page, scrapes each listed course page, and persists their lectures.
func (uc *scraperUsecase) ScrapeTopPageAndSave(ctx context.Context, year int) ([]domain.Lecture, error) {
	if uc.fetcher == nil {
		return nil, errors.New("scraper fetcher is not initialized")
	}

	reader, err := uc.fetcher.Fetch(ctx, scraper.TopPageURL)
	if err != nil {
		return nil, fmt.Errorf("fetch top page %s: %w", scraper.TopPageURL, err)
	}
	if reader == nil {
		return nil, fmt.Errorf("fetch top page %s: empty response", scraper.TopPageURL)
	}
	defer reader.Close()

	urls, err := uc.parser.ListCoursesPagesURL(reader, year)
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return []domain.Lecture{}, nil
	}

	seen := make(map[string]struct{}, len(urls))
	aggregated := make([]domain.Lecture, 0)
	firstList := true

	for _, listURL := range urls {
		if ctx != nil && ctx.Err() != nil {
			return nil, ctx.Err()
		}

		listURL = strings.TrimSpace(listURL)
		if listURL == "" {
			continue
		}
		if _, ok := seen[listURL]; ok {
			continue
		}
		seen[listURL] = struct{}{}

		if !firstList {
			if err := uc.sleep(ctx); err != nil {
				return nil, err
			}
		}
		firstList = false

		lectures, err := uc.ScrapeCourseListAndSave(ctx, listURL, scraper.TopPageURL)
		if err != nil {
			return nil, fmt.Errorf("scrape course list %s: %w", listURL, err)
		}
		aggregated = append(aggregated, lectures...)
	}

	return aggregated, nil
}

func (uc *scraperUsecase) sleep(ctx context.Context) error {
	if uc.delay <= 0 {
		return nil
	}

	timer := time.NewTimer(uc.delay)
	defer timer.Stop()

	if ctx == nil {
		<-timer.C
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
