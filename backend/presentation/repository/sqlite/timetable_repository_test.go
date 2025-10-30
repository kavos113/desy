package sqlite

import (
	"context"
	"testing"

	"github.com/kavos113/desy/backend/domain"
)

func TestTimetableRepositoryExpandTimetableRanges(t *testing.T) {
	_, db := newTestRepository(t)

	repo, err := NewTimetableRepository(db)
	if err != nil {
		t.Fatalf("NewTimetableRepository returned error: %v", err)
	}

	mustExec(t, db, `INSERT INTO lectures (id, university, title) VALUES (?, ?, ?)`, 1, "Test University", "Range Course")
	mustExec(t, db, `INSERT INTO rooms (id, name) VALUES (?, ?)`, 1, "Lecture Hall")

	mustExec(t, db, `INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`, 1, string(domain.SemesterFall), 1, string(domain.DayOfWeekMonday), 5)
	mustExec(t, db, `INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`, 1, string(domain.SemesterFall), 1, string(domain.DayOfWeekMonday), 8)

	inserted, err := repo.ExpandTimetableRanges(context.Background())
	if err != nil {
		t.Fatalf("ExpandTimetableRanges returned error: %v", err)
	}
	if inserted != 2 {
		t.Fatalf("expected 2 inserted periods, got %d", inserted)
	}

	rows, err := db.Query(`SELECT period FROM timetables WHERE lecture_id = ? ORDER BY period`, 1)
	if err != nil {
		t.Fatalf("select expanded timetables: %v", err)
	}
	defer rows.Close()

	var periods []int
	for rows.Next() {
		var period int
		if err := rows.Scan(&period); err != nil {
			rows.Close()
			t.Fatalf("scan expanded period: %v", err)
		}
		periods = append(periods, period)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate expanded periods: %v", err)
	}

	expected := []int{5, 6, 7, 8}
	if len(periods) != len(expected) {
		t.Fatalf("unexpected period count: got %d want %d", len(periods), len(expected))
	}
	for idx, period := range expected {
		if periods[idx] != period {
			t.Fatalf("unexpected period at %d: got %d want %d", idx, periods[idx], period)
		}
	}

	insertedAgain, err := repo.ExpandTimetableRanges(context.Background())
	if err != nil {
		t.Fatalf("second ExpandTimetableRanges returned error: %v", err)
	}
	if insertedAgain != 0 {
		t.Fatalf("expected idempotent second run, got %d new rows", insertedAgain)
	}
}

func TestTimetableRepositoryExpandTimetableRangesDifferentDay(t *testing.T) {
	_, db := newTestRepository(t)

	repo, err := NewTimetableRepository(db)
	if err != nil {
		t.Fatalf("NewTimetableRepository returned error: %v", err)
	}

	mustExec(t, db, `INSERT INTO lectures (id, university, title) VALUES (?, ?, ?)`, 1, "Test University", "Different Day")

	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 1, string(domain.DayOfWeekMonday), 5)
	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 1, string(domain.DayOfWeekTuesday), 8)

	inserted, err := repo.ExpandTimetableRanges(context.Background())
	if err != nil {
		t.Fatalf("ExpandTimetableRanges returned error: %v", err)
	}
	if inserted != 0 {
		t.Fatalf("expected no insertion for different days, got %d", inserted)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM timetables WHERE lecture_id = ?`, 1).Scan(&count); err != nil {
		t.Fatalf("count timetables: %v", err)
	}
	if count != 2 {
		t.Fatalf("unexpected timetable count: got %d want %d", count, 2)
	}
}

func TestTimetableRepositoryExpandTimetableRangesRequiresOddEven(t *testing.T) {
	_, db := newTestRepository(t)

	repo, err := NewTimetableRepository(db)
	if err != nil {
		t.Fatalf("NewTimetableRepository returned error: %v", err)
	}

	mustExec(t, db, `INSERT INTO lectures (id, university, title) VALUES (?, ?, ?)`, 1, "Test University", "Non Odd Even")
	mustExec(t, db, `INSERT INTO rooms (id, name) VALUES (?, ?)`, 1, "Hall")

	periods := []int{3, 4, 7, 8}
	for _, p := range periods {
		mustExec(t, db, `INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`, 1, string(domain.SemesterFall), 1, string(domain.DayOfWeekMonday), p)
	}

	inserted, err := repo.ExpandTimetableRanges(context.Background())
	if err != nil {
		t.Fatalf("ExpandTimetableRanges returned error: %v", err)
	}
	if inserted != 0 {
		t.Fatalf("expected no insertion for non odd-even range, got %d", inserted)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM timetables WHERE lecture_id = ?`, 1).Scan(&count); err != nil {
		t.Fatalf("count timetables: %v", err)
	}
	if count != len(periods) {
		t.Fatalf("unexpected timetable count: got %d want %d", count, len(periods))
	}
}

func TestTimetableRepositoryExpandTimetableRangesSpecialCaseTwoFour(t *testing.T) {
	_, db := newTestRepository(t)

	repo, err := NewTimetableRepository(db)
	if err != nil {
		t.Fatalf("NewTimetableRepository returned error: %v", err)
	}

	mustExec(t, db, `INSERT INTO lectures (id, university, title) VALUES (?, ?, ?)`, 1, "Test University", "Special Case")

	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 1, string(domain.DayOfWeekWednesday), 2)
	mustExec(t, db, `INSERT INTO timetables (lecture_id, day_of_week, period) VALUES (?, ?, ?)`, 1, string(domain.DayOfWeekWednesday), 4)

	inserted, err := repo.ExpandTimetableRanges(context.Background())
	if err != nil {
		t.Fatalf("ExpandTimetableRanges returned error: %v", err)
	}
	if inserted != 1 {
		t.Fatalf("expected single insertion for 2-4 special case, got %d", inserted)
	}

	rows, err := db.Query(`SELECT period FROM timetables WHERE lecture_id = ? ORDER BY period`, 1)
	if err != nil {
		t.Fatalf("select timetables: %v", err)
	}
	defer rows.Close()

	var periods []int
	for rows.Next() {
		var p int
		if err := rows.Scan(&p); err != nil {
			rows.Close()
			t.Fatalf("scan period: %v", err)
		}
		periods = append(periods, p)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate periods: %v", err)
	}

	expected := []int{2, 3, 4}
	if len(periods) != len(expected) {
		t.Fatalf("unexpected period count: got %d want %d", len(periods), len(expected))
	}
	for idx, period := range expected {
		if periods[idx] != period {
			t.Fatalf("unexpected period at %d: got %d want %d", idx, periods[idx], period)
		}
	}
}
