package sqlite

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kavos113/desy/backend/domain"
	"github.com/kavos113/desy/backend/presentation/scraper"
)

const testDataSourceName = "file::memory:?cache=shared"

func TestNewLectureRepositoryInitializesSchema(t *testing.T) {
	repo, db := newTestRepository(t)
	if repo == nil {
		t.Fatalf("expected repository instance")
	}

	tables := []string{
		"lectures",
		"teachers",
		"lecture_teachers",
		"rooms",
		"timetables",
		"lecture_plans",
		"lecture_keywords",
		"related_courses",
		"related_course_codes",
	}

	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Fatalf("table %s should exist: %v", table, err)
		}
		if name != table {
			t.Fatalf("unexpected table name for %s: %s", table, name)
		}
	}
}

func TestLectureRepositoryFindByIDReturnsAggregate(t *testing.T) {
	repo, db := newTestRepository(t)
	seedLectureAggregate(t, db)

	lecture, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if lecture == nil {
		t.Fatalf("expected lecture, got nil")
	}

	if lecture.Title != "Advanced Data Science" {
		t.Fatalf("unexpected title: %s", lecture.Title)
	}
	if lecture.EnglishTitle != "Advanced Data Science" {
		t.Fatalf("unexpected english title: %s", lecture.EnglishTitle)
	}
	if lecture.Department != "Computer Science" {
		t.Fatalf("unexpected department: %s", lecture.Department)
	}
	if lecture.LectureType != domain.LectureTypeLive {
		t.Fatalf("unexpected lecture type: %s", lecture.LectureType)
	}
	if lecture.Level != domain.LevelBachelor1 {
		t.Fatalf("unexpected level: %d", lecture.Level)
	}
	if lecture.Credit != 2 {
		t.Fatalf("unexpected credit: %d", lecture.Credit)
	}
	if lecture.Year != 2025 {
		t.Fatalf("unexpected year: %d", lecture.Year)
	}
	if lecture.Language != "English" {
		t.Fatalf("unexpected language: %s", lecture.Language)
	}
	if lecture.Url != "https://example.com/lectures/1" {
		t.Fatalf("unexpected url: %s", lecture.Url)
	}
	if lecture.Textbook != "Data Science Handbook" {
		t.Fatalf("unexpected textbook: %s", lecture.Textbook)
	}
	if lecture.Assessment != "Final Exam" {
		t.Fatalf("unexpected assessment: %s", lecture.Assessment)
	}
	if len(lecture.Timetables) != 1 {
		t.Fatalf("unexpected timetables length: %d", len(lecture.Timetables))
	}
	timetable := lecture.Timetables[0]
	if timetable.Semester != domain.SemesterSpring {
		t.Fatalf("unexpected semester: %s", timetable.Semester)
	}
	if timetable.DayOfWeek != domain.DayOfWeekMonday {
		t.Fatalf("unexpected day of week: %s", timetable.DayOfWeek)
	}
	if timetable.Period != domain.Period1 {
		t.Fatalf("unexpected period: %d", timetable.Period)
	}
	if timetable.Room.ID != 101 {
		t.Fatalf("unexpected room id: %d", timetable.Room.ID)
	}
	if timetable.Room.Name != "Room A" {
		t.Fatalf("unexpected room name: %s", timetable.Room.Name)
	}

	if len(lecture.Teachers) != 1 {
		t.Fatalf("unexpected teachers length: %d", len(lecture.Teachers))
	}
	if lecture.Teachers[0].Name != "Dr. Alice" {
		t.Fatalf("unexpected teacher name: %s", lecture.Teachers[0].Name)
	}

	if len(lecture.LecturePlans) != 1 {
		t.Fatalf("unexpected plans length: %d", len(lecture.LecturePlans))
	}
	if lecture.LecturePlans[0].Plan != "Introduction to advanced topics" {
		t.Fatalf("unexpected plan: %s", lecture.LecturePlans[0].Plan)
	}

	if len(lecture.Keywords) != 2 {
		t.Fatalf("unexpected keywords length: %d", len(lecture.Keywords))
	}
	if lecture.Keywords[0] != "data-science" || lecture.Keywords[1] != "machine-learning" {
		t.Fatalf("unexpected keywords: %#v", lecture.Keywords)
	}

	if len(lecture.RelatedCourses) != 1 || lecture.RelatedCourses[0] != 2 {
		t.Fatalf("unexpected related courses: %#v", lecture.RelatedCourses)
	}
	if len(lecture.RelatedCourseCodes) != 1 || lecture.RelatedCourseCodes[0] != "CS102" {
		t.Fatalf("unexpected related course codes: %#v", lecture.RelatedCourseCodes)
	}
	if lecture.OpenTerm != "2025 3Q" {
		t.Fatalf("unexpected open term: %s", lecture.OpenTerm)
	}
	expectedUpdated := time.Date(2025, time.March, 19, 0, 0, 0, 0, time.UTC)
	if !lecture.UpdatedAt.Equal(expectedUpdated) {
		t.Fatalf("unexpected updated_at: %v", lecture.UpdatedAt)
	}
}

func TestLectureRepositoryFindByIDNotFound(t *testing.T) {
	repo, _ := newTestRepository(t)
	lecture, err := repo.FindByID(42)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if lecture != nil {
		t.Fatalf("expected nil lecture, got %#v", lecture)
	}
}

func TestLectureRepositorySearchAppliesFilters(t *testing.T) {
	repo, db := newTestRepository(t)
	seedSearchData(t, db)

	result, err := repo.Search(domain.SearchQuery{
		TeacherName: "alice",
		Keywords:    []string{"science"},
		TimeTables: []domain.TimeTable{
			{DayOfWeek: domain.DayOfWeekMonday, Period: domain.Period1},
		},
		Title:       "data",
		Departments: []string{"Computer Science"},
		Year:        2025,
		Levels:      []domain.Level{domain.LevelBachelor1},
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected single result, got %d", len(result))
	}

	summary := result[0]
	if summary.ID != 1 {
		t.Fatalf("unexpected lecture id: %d", summary.ID)
	}
	if summary.Title != "Data Science Basics" {
		t.Fatalf("unexpected title: %s", summary.Title)
	}
	if summary.Level != domain.LevelBachelor1 {
		t.Fatalf("unexpected level: %d", summary.Level)
	}
	if summary.Year != 2025 {
		t.Fatalf("unexpected year: %d", summary.Year)
	}
	if len(summary.Timetables) != 1 {
		t.Fatalf("unexpected timetable length: %d", len(summary.Timetables))
	}
	if summary.Timetables[0].DayOfWeek != domain.DayOfWeekMonday {
		t.Fatalf("unexpected timetable day: %s", summary.Timetables[0].DayOfWeek)
	}
	if len(summary.Teachers) != 1 || summary.Teachers[0].Name != "Alice Smith" {
		t.Fatalf("unexpected teacher info: %#v", summary.Teachers)
	}
}

func TestLectureRepositorySearchReturnsEmptyWhenNoMatches(t *testing.T) {
	repo, _ := newTestRepository(t)
	result, err := repo.Search(domain.SearchQuery{Title: "non-existent"})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestLectureRepositoryCreatePersistsAggregate(t *testing.T) {
	repo, _ := newTestRepository(t)
	lecture := parseDetailFixture(t, "course_detail.html")

	if err := repo.Create(lecture); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if lecture.ID == 0 {
		t.Fatalf("expected ID to be assigned")
	}

	saved, err := repo.FindByID(lecture.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if saved == nil {
		t.Fatalf("expected saved lecture")
	}

	if saved.Title != lecture.Title {
		t.Fatalf("unexpected title: %s", saved.Title)
	}
	if len(saved.Teachers) != len(lecture.Teachers) {
		t.Fatalf("unexpected teachers count: %d", len(saved.Teachers))
	}
	if len(saved.Timetables) != len(lecture.Timetables) {
		t.Fatalf("unexpected timetables count: %d", len(saved.Timetables))
	}
	if len(saved.LecturePlans) == 0 {
		t.Fatalf("expected lecture plans to be stored")
	}
	if len(saved.Keywords) == 0 {
		t.Fatalf("expected keywords to be stored")
	}
	if len(saved.RelatedCourseCodes) == 0 {
		t.Fatalf("expected related course codes to be stored")
	}
	if saved.OpenTerm != strings.TrimSpace(lecture.OpenTerm) {
		t.Fatalf("unexpected open term: got %s want %s", saved.OpenTerm, lecture.OpenTerm)
	}
	if saved.UpdatedAt.IsZero() {
		t.Fatalf("expected updated_at to be stored")
	}
	if !saved.UpdatedAt.Equal(lecture.UpdatedAt) {
		t.Fatalf("unexpected updated_at: got %v want %v", saved.UpdatedAt, lecture.UpdatedAt)
	}
}

func TestLectureRepositoryCreateResolvesRelatedCourses(t *testing.T) {
	repo, db := newTestRepository(t)
	mustExec(t, db, `INSERT INTO lectures (id, university, title, code) VALUES (?, ?, ?, ?)`, 100, "Test University", "Existing Course", "LAH.S201")

	lecture := parseDetailFixture(t, "course_detail.html")

	if err := repo.Create(lecture); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if lecture.ID == 0 {
		t.Fatalf("expected lecture ID to be assigned")
	}

	saved, err := repo.FindByID(lecture.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if saved == nil {
		t.Fatalf("expected saved lecture")
	}
	if len(saved.RelatedCourses) == 0 {
		t.Fatalf("expected related courses to be resolved")
	}
	if saved.RelatedCourses[0] != 100 {
		t.Fatalf("unexpected related course id: %#v", saved.RelatedCourses)
	}
	if len(saved.RelatedCourseCodes) == 0 {
		t.Fatalf("expected related course codes to be stored")
	}
}

func TestLectureRepositoryMigrateRelatedCourses(t *testing.T) {
	repo, db := newTestRepository(t)

	mustExec(t, db, `INSERT INTO lectures (id, university, title, code) VALUES (?, ?, ?, ?)`, 1, "Test University", "Course A", "AAA100")
	mustExec(t, db, `INSERT INTO lectures (id, university, title, code) VALUES (?, ?, ?, ?)`, 2, "Test University", "Course B", "AAA200")
	mustExec(t, db, `INSERT INTO lectures (id, university, title, code) VALUES (?, ?, ?, ?)`, 3, "Test University", "Course C", "AAA300")

	mustExec(t, db, `INSERT INTO related_courses (lecture_id, related_lecture_id) VALUES (?, ?)`, 1, 2)

	codes := []struct {
		lectureID int
		code      string
	}{
		{lectureID: 1, code: "AAA200"},
		{lectureID: 1, code: "AAA300"},
		{lectureID: 2, code: "AAA300"},
		{lectureID: 3, code: "AAA999"}, // no matching lecture
		{lectureID: 3, code: "aaa100"},
	}

	for _, entry := range codes {
		mustExec(t, db, `INSERT INTO related_course_codes (lecture_id, code) VALUES (?, ?)`, entry.lectureID, entry.code)
	}

	inserted, err := repo.MigrateRelatedCourses(context.Background())
	if err != nil {
		t.Fatalf("MigrateRelatedCourses returned error: %v", err)
	}
	if inserted != 3 {
		t.Fatalf("expected 3 new relations, got %d", inserted)
	}

	rows, err := db.Query(`SELECT lecture_id, related_lecture_id FROM related_courses ORDER BY lecture_id, related_lecture_id`)
	if err != nil {
		t.Fatalf("select related courses: %v", err)
	}
	defer rows.Close()

	var relations [][2]int
	for rows.Next() {
		var src, dst int
		if err := rows.Scan(&src, &dst); err != nil {
			t.Fatalf("scan relation: %v", err)
		}
		relations = append(relations, [2]int{src, dst})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate relations: %v", err)
	}

	expected := [][2]int{{1, 2}, {1, 3}, {2, 3}, {3, 1}}
	if len(relations) != len(expected) {
		t.Fatalf("unexpected relation count: got %d want %d", len(relations), len(expected))
	}
	for idx, pair := range expected {
		if relations[idx] != pair {
			t.Fatalf("unexpected relation at %d: got %v want %v", idx, relations[idx], pair)
		}
	}
}

func newTestRepository(t *testing.T) (*LectureRepository, *sql.DB) {
	t.Helper()

	db, err := sql.Open(testDriverName, testDataSourceName)
	if err != nil {
		t.Fatalf("open in-memory sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	repo, err := NewLectureRepository(db)
	if err != nil {
		db.Close()
		t.Fatalf("NewLectureRepository: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return repo, db
}

func seedLectureAggregate(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExec(t, db, `INSERT INTO lectures (id, university, title, english_title, department, lecture_type, code, level, credit, year, language, url, abstract, goal, experience, flow, out_of_class_work, textbook, reference_book, assessment, prerequisite, contact, office_hours, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		1,
		"Test University",
		"Advanced Data Science",
		"Advanced Data Science",
		"Computer Science",
		string(domain.LectureTypeLive),
		"CS101",
		int(domain.LevelBachelor1),
		2,
		2025,
		"English",
		"https://example.com/lectures/1",
		"Course abstract",
		"Course goal",
		"Course experience",
		"Course flow",
		"Course work",
		"Data Science Handbook",
		"Reference",
		"Final Exam",
		"Prerequisites",
		"Contact Info",
		"Office Hours",
		"Notes",
	)
	mustExec(t, db, `UPDATE lectures SET open_term = ?, updated_at = ? WHERE id = ?`, "2025 3Q", "2025-03-19", 1)

	mustExec(t, db, `INSERT INTO lectures (id, university, title) VALUES (?, ?, ?)`, 2, "Test University", "Supporting Course")

	mustExec(t, db, `INSERT INTO teachers (id, name, url) VALUES (?, ?, ?)`, 1, "Dr. Alice", "https://example.com/teachers/1")
	mustExec(t, db, `INSERT INTO lecture_teachers (lecture_id, teacher_id) VALUES (?, ?)`, 1, 1)

	mustExec(t, db, `INSERT INTO rooms (id, name) VALUES (?, ?)`, 101, "Room A")

	mustExec(t, db, `INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`,
		1,
		string(domain.SemesterSpring),
		101,
		string(domain.DayOfWeekMonday),
		int(domain.Period1),
	)

	mustExec(t, db, `INSERT INTO lecture_plans (lecture_id, count, plan, assignment) VALUES (?, ?, ?, ?)`, 1, 1, "Introduction to advanced topics", "Read chapter 1")

	mustExec(t, db, `INSERT INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, 1, "machine-learning")
	mustExec(t, db, `INSERT INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, 1, "data-science")

	mustExec(t, db, `INSERT INTO related_courses (lecture_id, related_lecture_id) VALUES (?, ?)`, 1, 2)
	mustExec(t, db, `INSERT INTO related_course_codes (lecture_id, code) VALUES (?, ?)`, 1, "CS102")
}

func seedSearchData(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExec(t, db, `INSERT INTO lectures (id, university, title, english_title, department, code, level, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		1,
		"Test University",
		"Data Science Basics",
		"Introduction to Data",
		"Computer Science",
		"CS100",
		int(domain.LevelBachelor1),
		2025,
	)

	mustExec(t, db, `INSERT INTO lectures (id, university, title, department, code, level, year) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		2,
		"Test University",
		"Quantum Mechanics",
		"Physics",
		"PH200",
		int(domain.LevelBachelor2),
		2024,
	)

	mustExec(t, db, `INSERT INTO teachers (id, name, url) VALUES (?, ?, ?)`, 1, "Alice Smith", "https://example.com/teachers/alice")
	mustExec(t, db, `INSERT INTO teachers (id, name, url) VALUES (?, ?, ?)`, 2, "Bob Brown", "https://example.com/teachers/bob")

	mustExec(t, db, `INSERT INTO lecture_teachers (lecture_id, teacher_id) VALUES (?, ?)`, 1, 1)
	mustExec(t, db, `INSERT INTO lecture_teachers (lecture_id, teacher_id) VALUES (?, ?)`, 2, 2)

	mustExec(t, db, `INSERT INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, 1, "science")
	mustExec(t, db, `INSERT INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, 1, "data")
	mustExec(t, db, `INSERT INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, 2, "physics")

	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 1, string(domain.DayOfWeekMonday), int(domain.Period1))
	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 2, string(domain.DayOfWeekTuesday), int(domain.Period3))
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

func parseDetailFixture(t *testing.T, name string) *domain.Lecture {
	t.Helper()
	path := filepath.Join("..", "..", "scraper", "fixture", name)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture %s: %v", name, err)
	}
	defer file.Close()

	lecture, err := scraper.ParseCourseDetail(file, "https://example.com/courses/2025/LAH.S101")
	if err != nil {
		t.Fatalf("parse course detail: %v", err)
	}

	return lecture
}
