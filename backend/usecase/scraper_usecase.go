package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
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

// ScrapeProgress describes the current scraping status.
type ScrapeProgress struct {
	Total   int
	Current int
	Code    string
	Title   string
}

// ScrapeProgressReporter receives updates while scraping progresses.
type ScrapeProgressReporter interface {
	Report(ScrapeProgress)
}

// ScraperUsecase orchestrates scraping workflow and persistence.
type ScraperUsecase interface {
	ScrapeCourseList(ctx context.Context, listURL, baseURL string) ([]scraper.CourseListItem, error)
	ScrapeCourseListAndSave(ctx context.Context, listURL, baseURL string) ([]domain.Lecture, error)
	ScrapeCourseDetail(ctx context.Context, detailURL string) (*domain.Lecture, error)
	ScrapeCourseDetailAndSave(ctx context.Context, detailURL string) (*domain.Lecture, error)
	ScrapeTopPageAndSave(ctx context.Context, year int) ([]domain.Lecture, error)
	SetProgressReporter(ScrapeProgressReporter)
}

type scraperUsecase struct {
	fetcher     Fetcher
	lectureRepo domain.LectureRepository
	parser      scraper.Parser
	delay       time.Duration
	reporter    ScrapeProgressReporter
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

// SetProgressReporter registers a progress reporter to receive scraping updates.
func (uc *scraperUsecase) SetProgressReporter(reporter ScrapeProgressReporter) {
	uc.reporter = reporter
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
		uc.reportProgress(ScrapeProgress{Total: 0})
		return []domain.Lecture{}, nil
	}

	uniqueItems := make([]scraper.CourseListItem, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		detailURL := strings.TrimSpace(item.DetailURL)
		if detailURL == "" {
			continue
		}
		if _, ok := seen[detailURL]; ok {
			continue
		}
		seen[detailURL] = struct{}{}
		item.DetailURL = detailURL
		uniqueItems = append(uniqueItems, item)
	}

	total := len(uniqueItems)
	if total == 0 {
		uc.reportProgress(ScrapeProgress{Total: 0})
		return []domain.Lecture{}, nil
	}

	uc.reportProgress(ScrapeProgress{Total: total})

	lectures := make([]domain.Lecture, 0, total)
	firstFetch := true

	for idx, item := range uniqueItems {
		if ctx != nil && ctx.Err() != nil {
			return nil, ctx.Err()
		}

		detailURL := item.DetailURL
		log.Printf("Scraping detail page: %s %s", item.Code, item.Title)
		uc.reportProgress(ScrapeProgress{Total: total, Current: idx + 1, Code: strings.TrimSpace(item.Code), Title: strings.TrimSpace(item.Title)})

		existing, err := uc.lectureRepo.FindByCode(item.Code, item.Title, item.OpenTerm)
		if err != nil {
			return nil, fmt.Errorf("find lecture by code %s: %w", item.Code, err)
		}

		if shouldSkipLecture(existing, item) {
			continue
		}

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
		if openTerm := strings.TrimSpace(item.OpenTerm); openTerm != "" {
			lecture.OpenTerm = openTerm
		}
		lecture.UpdatedAt = normalizeDate(selectUpdatedAt(lecture.UpdatedAt, item.UpdatedAt))
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

	time.Sleep(1 * time.Second)

	// english title
	englishURL := buildEnglishDetailURL(detailURL)
	if englishURL != "" {
		engReader, err := uc.fetcher.Fetch(ctx, englishURL)
		if err != nil {
			log.Printf("fetch english title %s: %v", englishURL, err)
		} else {
			defer engReader.Close()
			if err := uc.parser.AddEnglishTitle(engReader, lecture); err != nil {
				log.Printf("add english title from %s: %v", englishURL, err)
			}
		}
	}

	lecture.UpdatedAt = normalizeDate(lecture.UpdatedAt)

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
	lecture.UpdatedAt = normalizeDate(lecture.UpdatedAt)

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

	if _, err := uc.lectureRepo.MigrateRelatedCourses(ctx); err != nil {
		return nil, fmt.Errorf("migrate related courses: %w", err)
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

func buildEnglishDetailURL(detailURL string) string {
	detailURL = strings.TrimSpace(detailURL)
	if detailURL == "" {
		return ""
	}
	u, err := url.Parse(detailURL)
	if err != nil {
		if strings.Contains(detailURL, "?") {
			return detailURL + "&hl=en"
		}
		return detailURL + "?hl=en"
	}
	query := u.Query()
	query.Set("hl", "en")
	u.RawQuery = query.Encode()
	return u.String()
}

func (uc *scraperUsecase) reportProgress(progress ScrapeProgress) {
	if uc.reporter == nil {
		return
	}
	uc.reporter.Report(progress)
}

func shouldSkipLecture(existing *domain.Lecture, item scraper.CourseListItem) bool {
	if existing == nil {
		return false
	}

	if normalizeComparable(existing.Title) != normalizeComparable(item.Title) {
		return false
	}
	if normalizeComparable(existing.Code) != normalizeComparable(item.Code) {
		return false
	}
	if normalizeComparable(existing.OpenTerm) != normalizeComparable(item.OpenTerm) {
		return false
	}
	if !sameDate(existing.UpdatedAt, item.UpdatedAt) {
		return false
	}

	return true
}

func normalizeComparable(value string) string {
	if value == "" {
		return ""
	}
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}

func sameDate(a, b time.Time) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	if a.IsZero() || b.IsZero() {
		return false
	}
	return normalizeDate(a).Equal(normalizeDate(b))
}

func normalizeDate(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func selectUpdatedAt(detail, list time.Time) time.Time {
	if !list.IsZero() {
		return list
	}
	return detail
}
