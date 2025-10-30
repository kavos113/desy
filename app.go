package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/repository/sqlite"
	"github.com/kavos113/desy/backend/presentation/scraper"
	"github.com/kavos113/desy/backend/usecase"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx              context.Context
	db               *sql.DB
	lectureUsecase   usecase.LectureUsecase
	scraperUsecase   usecase.ScraperUsecase
	timetableUsecase usecase.TimeTableUsecase
}

// NewApp creates a new App application struct
func NewApp() *App {
	db, err := sql.Open("sqlite", "file:desy.db?_pragma=foreign_keys(1)")
	if err != nil {
		panic(fmt.Errorf("open sqlite database: %w", err))
	}

	lectureRepo, err := sqlite.NewLectureRepository(db)
	if err != nil {
		panic(fmt.Errorf("init lecture repository: %w", err))
	}

	timetableRepo, err := sqlite.NewTimetableRepository(db)
	if err != nil {
		panic(fmt.Errorf("init timetable repository: %w", err))
	}

	fetcher := usecase.NewHTTPFetcher(&http.Client{Timeout: 15 * time.Second})
	scraperUsecase := usecase.NewScraperUsecase(fetcher, lectureRepo, timetableRepo, scraper.NewParser(), 3*time.Second)

	return &App{
		db:               db,
		lectureUsecase:   usecase.NewLectureUsecase(lectureRepo),
		scraperUsecase:   scraperUsecase,
		timetableUsecase: usecase.NewTimeTableUsecase(timetableRepo),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if a.timetableUsecase != nil {
		if _, err := a.timetableUsecase.ExpandTimetableRanges(ctx); err != nil {
			log.Printf("expand timetable ranges on startup: %v", err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Scrape() error {
	if a.scraperUsecase == nil {
		return fmt.Errorf("scraper usecase is not configured")
	}

	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cleanup := a.attachProgressReporter(ctx)
	defer cleanup()

	_, err := a.scraperUsecase.ScrapeTopPageAndSave(ctx, time.Now().Year())
	return err
}

func (a *App) ScrapeAll() error {
	if a.scraperUsecase == nil {
		return fmt.Errorf("scraper usecase is not configured")
	}

	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cleanup := a.attachProgressReporter(ctx)
	defer cleanup()

	for year := 2020; year <= time.Now().Year(); year++ {
		_, err := a.scraperUsecase.ScrapeTopPageAndSave(ctx, year)
		if err != nil {
			return fmt.Errorf("scrape year %d: %w", year, err)
		}
	}
	return nil
}

func (a *App) ScrapeTest() error {
	const testURL = "https://syllabus.s.isct.ac.jp/courses/2025/4/0-904-340000-120900-20927"
	if a.scraperUsecase == nil {
		return fmt.Errorf("scraper usecase is not configured")
	}

	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cleanup := a.attachProgressReporter(ctx)
	defer cleanup()

	_, err := a.scraperUsecase.ScrapeCourseListAndSave(ctx, testURL, scraper.TopPageURL)
	return err
}

func (a *App) MigrateRelatedCourses() (int, error) {
	if a.lectureUsecase == nil {
		return 0, fmt.Errorf("lecture usecase is not configured")
	}

	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	return a.lectureUsecase.MigrateRelatedCourses(ctx)
}

func (a *App) SearchLectures(query domain.SearchQuery) ([]domain.LectureSummary, error) {
	if a.lectureUsecase == nil {
		return nil, fmt.Errorf("lecture usecase is not configured")
	}

	return a.lectureUsecase.SearchLectures(query)
}

func (a *App) GetLectureDetails(lectureID int) (*domain.Lecture, error) {
	if a.lectureUsecase == nil {
		return nil, fmt.Errorf("lecture usecase is not configured")
	}

	return a.lectureUsecase.GetLectureDetails(lectureID)
}

func (a *App) shutdown(context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func (a *App) attachProgressReporter(ctx context.Context) func() {
	if a.scraperUsecase == nil || ctx == nil {
		return func() {}
	}
	reporter := &wailsProgressReporter{ctx: ctx}
	a.scraperUsecase.SetProgressReporter(reporter)
	return func() {
		a.scraperUsecase.SetProgressReporter(nil)
	}
}

type wailsProgressReporter struct {
	ctx context.Context
}

func (r *wailsProgressReporter) Report(progress usecase.ScrapeProgress) {
	if r == nil || r.ctx == nil {
		return
	}
	runtime.EventsEmit(r.ctx, "fetch_status", progress)
}
