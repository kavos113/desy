package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/repository/sqlite"
	"github.com/kavos113/desy/backend/presentation/scraper"
	_ "modernc.org/sqlite"
)

const testDataSourceName = "file::memory:?cache=shared"

func TestScraperUsecaseScrapeCourseDetailAndSave(t *testing.T) {
	repo, _ := newUsecaseTestRepository(t)
	detailURL := "https://example.com/courses/2025/LAH.S101"

	fetcher := newMockFetcher(map[string]string{
		detailURL: readFixture(t, "course_detail.html"),
	})

	usecase := NewScraperUsecase(fetcher, repo, scraper.NewParser())

	lecture, err := usecase.ScrapeCourseDetailAndSave(context.Background(), detailURL)
	if err != nil {
		t.Fatalf("ScrapeCourseDetailAndSave returned error: %v", err)
	}
	if lecture == nil {
		t.Fatalf("expected lecture, got nil")
	}
	if lecture.ID == 0 {
		t.Fatalf("expected lecture ID to be assigned")
	}
	if len(lecture.Teachers) == 0 {
		t.Fatalf("expected teachers to be populated")
	}

	stored, err := repo.FindByID(lecture.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if stored == nil {
		t.Fatalf("expected stored lecture")
	}
	if stored.Title != lecture.Title {
		t.Fatalf("unexpected stored title: %s", stored.Title)
	}
	if len(stored.Timetables) == 0 {
		t.Fatalf("expected stored timetables")
	}
}

func TestScraperUsecaseScrapeCourseListAndSave(t *testing.T) {
	repo, _ := newUsecaseTestRepository(t)

	listURL := "https://example.com/list"
	detailURL := "https://example.com/courses/2025/LAH.S101"
	baseURL := "https://example.com"

	listHTML := `
<html><body>
<table class="c-table">
  <tbody>
    <tr>
      <td>LAH.S101</td>
      <td><a href="/courses/2025/LAH.S101">法学（憲法）Ａ</a></td>
    </tr>
  </tbody>
</table>
</body></html>`

	fetcher := newMockFetcher(map[string]string{
		listURL:   listHTML,
		detailURL: readFixture(t, "course_detail.html"),
	})

	usecase := NewScraperUsecase(fetcher, repo, scraper.NewParser())

	lectures, err := usecase.ScrapeCourseListAndSave(context.Background(), listURL, baseURL)
	if err != nil {
		t.Fatalf("ScrapeCourseListAndSave returned error: %v", err)
	}
	if len(lectures) != 1 {
		t.Fatalf("expected single lecture, got %d", len(lectures))
	}
	if lectures[0].ID == 0 {
		t.Fatalf("expected lecture ID to be assigned")
	}

	stored, err := repo.FindByID(lectures[0].ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if stored == nil {
		t.Fatalf("expected stored lecture")
	}
	if stored.Code != "LAH.S101" {
		t.Fatalf("unexpected stored code: %s", stored.Code)
	}
	if len(stored.Teachers) == 0 {
		t.Fatalf("expected stored teachers")
	}
}

func newUsecaseTestRepository(t *testing.T) (domain.LectureRepository, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", testDataSourceName)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	repo, err := sqlite.NewLectureRepository(db)
	if err != nil {
		db.Close()
		t.Fatalf("create lecture repository: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return repo, db
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("..", "presentation", "scraper", "fixture", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

type mockFetcher struct {
	responses map[string]string
}

func newMockFetcher(responses map[string]string) *mockFetcher {
	return &mockFetcher{responses: responses}
}

func (m *mockFetcher) Fetch(_ context.Context, url string) (io.ReadCloser, error) {
	body, ok := m.responses[url]
	if !ok {
		return nil, fmt.Errorf("unexpected fetch URL: %s", url)
	}
	return io.NopCloser(strings.NewReader(body)), nil
}
