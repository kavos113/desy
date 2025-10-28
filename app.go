package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	domain "github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/repository/sqlite"
	"github.com/kavos113/desy/backend/presentation/scraper"
	"github.com/kavos113/desy/backend/usecase"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx            context.Context
	db             *sql.DB
	lectureUsecase usecase.LectureUsecase
	scraperUsecase usecase.ScraperUsecase
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

	fetcher := usecase.NewHTTPFetcher(&http.Client{Timeout: 15 * time.Second})
	scraperUsecase := usecase.NewScraperUsecase(fetcher, lectureRepo, scraper.NewParser(), 3*time.Second)

	return &App{
		db:             db,
		lectureUsecase: usecase.NewLectureUsecase(lectureRepo),
		scraperUsecase: scraperUsecase,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	_, err := a.scraperUsecase.ScrapeTopPageAndSave(ctx, time.Now().Year())
	return err
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

	_, err := a.scraperUsecase.ScrapeCourseListAndSave(ctx, testURL, scraper.TopPageURL)
	return err
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
