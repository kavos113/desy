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
	"time"

	"github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/repository/sqlite"
	"github.com/kavos113/desy/backend/presentation/scraper"
	_ "modernc.org/sqlite"
)

const testDataSourceName = "file::memory:?cache=shared"

func TestScraperUsecaseScrapeCourseDetailAndSave(t *testing.T) {
	repo, timetableRepo, _ := newUsecaseTestRepository(t)
	detailURL := "https://example.com/courses/2025/LAH.S101"
	detailURLEnglish := buildEnglishDetailURL(detailURL)

	fetcher := newMockFetcher(map[string]string{
		detailURL:        readFixture(t, "course_detail.html"),
		detailURLEnglish: readFixture(t, "course_detail_en.html"),
	})

	usecase := NewScraperUsecase(fetcher, repo, timetableRepo, scraper.NewParser(), 0)

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
	if lecture.EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected english title: %s", lecture.EnglishTitle)
	}
	if lecture.OpenTerm != "2025 3Q" {
		t.Fatalf("unexpected open term: %s", lecture.OpenTerm)
	}
	expectedUpdated := time.Date(2025, time.March, 19, 0, 0, 0, 0, time.UTC)
	if !lecture.UpdatedAt.Equal(expectedUpdated) {
		t.Fatalf("unexpected updated_at: %v", lecture.UpdatedAt)
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
	if stored.EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected stored english title: %s", stored.EnglishTitle)
	}
	if stored.OpenTerm != "2025 3Q" {
		t.Fatalf("unexpected stored open term: %s", stored.OpenTerm)
	}
	if !stored.UpdatedAt.Equal(expectedUpdated) {
		t.Fatalf("unexpected stored updated_at: %v", stored.UpdatedAt)
	}
}

func TestScraperUsecaseScrapeCourseListAndSave(t *testing.T) {
	repo, timetableRepo, _ := newUsecaseTestRepository(t)

	listURL := "https://example.com/list"
	detailURL := "https://example.com/courses/2025/LAH.S101"
	detailURLEnglish := buildEnglishDetailURL(detailURL)
	baseURL := "https://example.com"

	listHTML := `
<html><body>
<table class="c-table">
  <tbody>
    <tr>
      <td>LAH.S101</td>
      <td><a href="/courses/2025/LAH.S101">法学（憲法）Ａ</a></td>
				<td>篠島 正幸</td>
				<td>文系教養科目</td>
				<td>2025 3Q</td>
				<td>2025/3/19</td>
    </tr>
  </tbody>
</table>
</body></html>`

	fetcher := newMockFetcher(map[string]string{
		listURL:          listHTML,
		detailURL:        readFixture(t, "course_detail.html"),
		detailURLEnglish: readFixture(t, "course_detail_en.html"),
	})

	usecase := NewScraperUsecase(fetcher, repo, timetableRepo, scraper.NewParser(), 0)
	reporter := &collectingProgressReporter{}
	usecase.SetProgressReporter(reporter)
	t.Cleanup(func() {
		usecase.SetProgressReporter(nil)
	})

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
	if lectures[0].EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected english title: %s", lectures[0].EnglishTitle)
	}
	if lectures[0].OpenTerm != "2025 3Q" {
		t.Fatalf("unexpected open term: %s", lectures[0].OpenTerm)
	}
	if lectures[0].UpdatedAt.IsZero() {
		t.Fatalf("expected updated_at to be captured")
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
	if stored.EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected stored english title: %s", stored.EnglishTitle)
	}
	if stored.OpenTerm != "2025 3Q" {
		t.Fatalf("unexpected stored open term: %s", stored.OpenTerm)
	}
	if stored.UpdatedAt.IsZero() {
		t.Fatalf("expected stored updated_at")
	}
	if !stored.UpdatedAt.Equal(lectures[0].UpdatedAt) {
		t.Fatalf("unexpected stored updated_at: %v", stored.UpdatedAt)
	}
	if len(reporter.events) != 2 {
		t.Fatalf("unexpected number of progress events: %d", len(reporter.events))
	}
	initial := reporter.events[0]
	if initial.Total != 1 {
		t.Fatalf("unexpected initial total: %+v", initial)
	}
	if initial.Current != 0 {
		t.Fatalf("unexpected initial current: %+v", initial)
	}
	if initial.Code != "" || initial.Title != "" {
		t.Fatalf("initial progress should not include code or title: %+v", initial)
	}
	update := reporter.events[1]
	if update.Total != 1 || update.Current != 1 {
		t.Fatalf("unexpected progress update: %+v", update)
	}
	if update.Code != "LAH.S101" {
		t.Fatalf("unexpected progress code: %s", update.Code)
	}
	if update.Title != "法学（憲法）Ａ" {
		t.Fatalf("unexpected progress title: %s", update.Title)
	}
}

func TestScraperUsecaseScrapeCourseListAndSaveSkipsUnchanged(t *testing.T) {
	repo, timetableRepo, _ := newUsecaseTestRepository(t)

	parser := scraper.NewParser()
	lecture, err := parser.ParseCourseDetail(strings.NewReader(readFixture(t, "course_detail.html")), "https://example.com/courses/2025/LAH.S101")
	if err != nil {
		t.Fatalf("parse detail fixture: %v", err)
	}
	if err := repo.Create(lecture); err != nil {
		t.Fatalf("seed lecture: %v", err)
	}

	listURL := "https://example.com/list"
	baseURL := "https://example.com"

	listHTML := `
<html><body>
<table class="c-table">
  <tbody>
    <tr>
      <td>LAH.S101</td>
      <td><a href="/courses/2025/LAH.S101">法学（憲法）Ａ</a></td>
      <td>篠島 正幸</td>
      <td>文系教養科目</td>
      <td>2025 3Q</td>
      <td>2025/3/19</td>
    </tr>
  </tbody>
</table>
</body></html>`

	fetcher := newMockFetcher(map[string]string{
		listURL: listHTML,
	})

	usecase := NewScraperUsecase(fetcher, repo, timetableRepo, scraper.NewParser(), 0)

	lectures, err := usecase.ScrapeCourseListAndSave(context.Background(), listURL, baseURL)
	if err != nil {
		t.Fatalf("ScrapeCourseListAndSave returned error: %v", err)
	}
	if len(lectures) != 0 {
		t.Fatalf("expected no new lectures, got %d", len(lectures))
	}
}

func TestScraperUsecaseScrapeTopPageAndSave(t *testing.T) {
	repo, timetableRepo, _ := newUsecaseTestRepository(t)

	listURL := scraper.TopPageURL + "/courses/2025/4/mock-list"
	detailURL := scraper.TopPageURL + "/courses/2025/LAH.S101"
	detailURLEnglish := buildEnglishDetailURL(detailURL)

	topPage := fmt.Sprintf(`
<html><body>
  <a href="%s">ListA</a>
  <a href="%s">ListA Duplicate</a>
</body></html>`, listURL, listURL)

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
		scraper.TopPageURL: topPage,
		listURL:            listHTML,
		detailURL:          readFixture(t, "course_detail.html"),
		detailURLEnglish:   readFixture(t, "course_detail_en.html"),
	})

	usecase := NewScraperUsecase(fetcher, repo, timetableRepo, scraper.NewParser(), 0)
	reporter := &collectingProgressReporter{}
	usecase.SetProgressReporter(reporter)
	t.Cleanup(func() {
		usecase.SetProgressReporter(nil)
	})

	lectures, err := usecase.ScrapeTopPageAndSave(context.Background(), 2025)
	if err != nil {
		t.Fatalf("ScrapeTopPageAndSave returned error: %v", err)
	}
	if len(lectures) != 1 {
		t.Fatalf("expected single lecture, got %d", len(lectures))
	}
	if lectures[0].EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected english title: %s", lectures[0].EnglishTitle)
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
	if stored.EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected stored english title: %s", stored.EnglishTitle)
	}
	if len(reporter.events) != 2 {
		t.Fatalf("unexpected number of progress events: %d", len(reporter.events))
	}
}

func newUsecaseTestRepository(t *testing.T) (domain.LectureRepository, domain.TimeTableRepository, *sql.DB) {
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

	timetableRepo, err := sqlite.NewTimetableRepository(db)
	if err != nil {
		db.Close()
		t.Fatalf("create timetable repository: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return repo, timetableRepo, db
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

type collectingProgressReporter struct {
	events []ScrapeProgress
}

func (c *collectingProgressReporter) Report(progress ScrapeProgress) {
	c.events = append(c.events, progress)
}
