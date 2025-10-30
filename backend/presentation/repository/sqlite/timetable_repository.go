package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/kavos113/desy/backend/domain"
)

// ErrTimetableNotImplemented indicates that the operation is not yet supported.
var ErrTimetableNotImplemented = errors.New("not implemented")

// TimetableRepository provides SQLite backed access to timetables.
type TimetableRepository struct {
	db *sql.DB
}

// NewTimetableRepository creates a timetable repository for the provided database handle.
func NewTimetableRepository(db *sql.DB) (*TimetableRepository, error) {
	if db == nil {
		return nil, errors.New("nil database handle")
	}
	return &TimetableRepository{db: db}, nil
}

// FindByLectureID retrieves timetables belonging to a lecture.
func (r *TimetableRepository) FindByLectureID(lectureID int) ([]domain.TimeTable, error) {
	if lectureID <= 0 {
		return nil, fmt.Errorf("invalid lecture id: %d", lectureID)
	}

	rows, err := r.db.Query(`SELECT tt.lecture_id, tt.semester, tt.room_id, r.name, tt.day_of_week, tt.period FROM timetables tt LEFT JOIN rooms r ON r.id = tt.room_id WHERE tt.lecture_id = ? ORDER BY tt.semester, tt.day_of_week, tt.period`, lectureID)
	if err != nil {
		return nil, fmt.Errorf("select timetables: %w", err)
	}
	defer rows.Close()

	timetables := make([]domain.TimeTable, 0)
	for rows.Next() {
		var (
			lectureIDValue int
			semesterValue  sql.NullString
			roomIDValue    sql.NullInt64
			roomNameValue  sql.NullString
			dayValue       sql.NullString
			periodValue    sql.NullInt64
		)
		if err := rows.Scan(&lectureIDValue, &semesterValue, &roomIDValue, &roomNameValue, &dayValue, &periodValue); err != nil {
			return nil, fmt.Errorf("scan timetable: %w", err)
		}

		timetable := domain.TimeTable{LectureID: lectureIDValue}
		if semesterValue.Valid {
			timetable.Semester = domain.Semester(strings.TrimSpace(semesterValue.String))
		}
		if dayValue.Valid {
			timetable.DayOfWeek = domain.DayOfWeek(strings.TrimSpace(dayValue.String))
		}
		if periodValue.Valid {
			timetable.Period = domain.Period(periodValue.Int64)
		}
		if roomIDValue.Valid {
			timetable.Room.ID = int(roomIDValue.Int64)
		}
		if roomNameValue.Valid {
			timetable.Room.Name = roomNameValue.String
		}
		timetables = append(timetables, timetable)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate timetables: %w", err)
	}

	return timetables, nil
}

// Create inserts a single timetable record.
func (r *TimetableRepository) Create(timetable *domain.TimeTable) error {
	if timetable == nil {
		return errors.New("nil timetable")
	}
	return r.Creates([]domain.TimeTable{*timetable})
}

// Creates inserts multiple timetables in a single transaction.
func (r *TimetableRepository) Creates(timetables []domain.TimeTable) error {
	if len(timetables) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin timetable transaction: %w", err)
	}

	for idx := range timetables {
		if timetables[idx].LectureID <= 0 {
			tx.Rollback()
			return fmt.Errorf("timetable lecture id must be positive")
		}

		roomID := 0
		if name := strings.TrimSpace(timetables[idx].Room.Name); name != "" {
			id, err := r.ensureRoomTx(tx, name)
			if err != nil {
				tx.Rollback()
				return err
			}
			roomID = id
			timetables[idx].Room.ID = id
		}

		if _, err := tx.Exec(`INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`,
			timetables[idx].LectureID,
			nullString(string(timetables[idx].Semester)),
			nullInt(roomID),
			nullString(string(timetables[idx].DayOfWeek)),
			nullInt(int(timetables[idx].Period)),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert timetable: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit timetable transaction: %w", err)
	}

	return nil
}

// Update modifies a timetable record.
func (r *TimetableRepository) Update(_ *domain.TimeTable) error {
	return ErrTimetableNotImplemented
}

// Delete removes timetable records for the specified lecture.
func (r *TimetableRepository) Delete(lectureID int) error {
	if lectureID <= 0 {
		return fmt.Errorf("invalid lecture id: %d", lectureID)
	}
	if _, err := r.db.Exec(`DELETE FROM timetables WHERE lecture_id = ?`, lectureID); err != nil {
		return fmt.Errorf("delete timetables: %w", err)
	}
	return nil
}

// ExpandTimetableRanges fills gaps between range endpoints stored as separate periods.
func (r *TimetableRepository) ExpandTimetableRanges(ctx context.Context) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin expand timetable ranges transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT lecture_id, semester, room_id, day_of_week, period FROM timetables WHERE period IS NOT NULL AND period > 0 ORDER BY lecture_id, semester, day_of_week, room_id, period`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("select timetables for expansion: %w", err)
	}
	defer rows.Close()

	type groupKey struct {
		lectureID     int
		semester      string
		semesterValid bool
		dayOfWeek     string
		dayValid      bool
		roomID        int
		roomValid     bool
	}

	type groupData struct {
		periods map[int]struct{}
	}

	groups := make(map[groupKey]*groupData)

	for rows.Next() {
		var (
			lectureIDValue int
			semesterValue  sql.NullString
			roomIDValue    sql.NullInt64
			dayValue       sql.NullString
			periodValue    sql.NullInt64
		)

		if err := rows.Scan(&lectureIDValue, &semesterValue, &roomIDValue, &dayValue, &periodValue); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("scan timetable for expansion: %w", err)
		}

		if !periodValue.Valid {
			continue
		}

		key := groupKey{
			lectureID:     lectureIDValue,
			semester:      strings.TrimSpace(semesterValue.String),
			semesterValid: semesterValue.Valid && strings.TrimSpace(semesterValue.String) != "",
			dayOfWeek:     strings.TrimSpace(dayValue.String),
			dayValid:      dayValue.Valid && strings.TrimSpace(dayValue.String) != "",
			roomID:        int(roomIDValue.Int64),
			roomValid:     roomIDValue.Valid && roomIDValue.Int64 > 0,
		}

		data, exists := groups[key]
		if !exists {
			data = &groupData{periods: make(map[int]struct{})}
			groups[key] = data
		}

		if periodValue.Int64 <= 0 {
			continue
		}
		data.periods[int(periodValue.Int64)] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("iterate timetables for expansion: %w", err)
	}

	if len(groups) == 0 {
		if err := tx.Commit(); err != nil {
			return 0, fmt.Errorf("commit expand timetable ranges transaction: %w", err)
		}
		return 0, nil
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("prepare insert expanded timetable: %w", err)
	}
	defer stmt.Close()

	inserted := 0

	for key, data := range groups {
		if len(data.periods) != 2 {
			continue
		}

		periods := make([]int, 0, len(data.periods))
		for period := range data.periods {
			periods = append(periods, period)
		}
		sort.Ints(periods)

		start := periods[0]
		end := periods[1]
		if end <= start+1 {
			continue
		}

		isOddEven := start%2 == 1 && end%2 == 0
		isSpecial := start == 2 && end == 4
		if !isOddEven && !isSpecial {
			continue
		}

		var semesterArg interface{}
		if key.semesterValid {
			semesterArg = key.semester
		}

		var dayArg interface{}
		if key.dayValid {
			dayArg = key.dayOfWeek
		}

		var roomArg interface{}
		if key.roomValid {
			roomArg = key.roomID
		}

		for period := start + 1; period < end; period++ {
			result, err := stmt.Exec(key.lectureID, semesterArg, roomArg, dayArg, period)
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("insert expanded timetable: %w", err)
			}
			affected, err := result.RowsAffected()
			if err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("rows affected for expanded timetable: %w", err)
			}
			inserted += int(affected)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit expand timetable ranges transaction: %w", err)
	}

	return inserted, nil
}

func (r *TimetableRepository) ensureRoomTx(tx *sql.Tx, name string) (int, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, nil
	}

	var id int
	err := tx.QueryRow(`SELECT id FROM rooms WHERE name = ?`, name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("select room: %w", err)
	}

	result, err := tx.Exec(`INSERT OR IGNORE INTO rooms (name) VALUES (?)`, name)
	if err != nil {
		return 0, fmt.Errorf("insert room: %w", err)
	}

	id64, err := result.LastInsertId()
	if err == nil && id64 != 0 {
		return int(id64), nil
	}

	err = tx.QueryRow(`SELECT id FROM rooms WHERE name = ?`, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("select ensured room: %w", err)
	}

	return id, nil
}
