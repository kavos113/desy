package main

import (
	"context"
	"database/sql"
	"fmt"

	domain "github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/repository/sqlite"
	"github.com/kavos113/desy/backend/usecase"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx            context.Context
	db             *sql.DB
	lectureUsecase usecase.LectureUsecase
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

	return &App{
		db:             db,
		lectureUsecase: usecase.NewLectureUsecase(lectureRepo),
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

func (a *App) Scrape() {

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
