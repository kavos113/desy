package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/kavos113/desy/backend/domain"
)

// ErrNotImplemented indicates that the repository operation is not yet implemented.
var ErrNotImplemented = errors.New("not implemented")

// LectureRepository provides SQLite backed access to lecture aggregates.
type LectureRepository struct {
	db *sql.DB
}

// NewLectureRepository creates a repository instance and ensures the required schema exists.
func NewLectureRepository(db *sql.DB) (*LectureRepository, error) {
	if db == nil {
		return nil, errors.New("nil database handle")
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	repo := &LectureRepository{db: db}
	if err := repo.initSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

// FindByID retrieves a lecture aggregate by identifier.
func (r *LectureRepository) FindByID(id int) (*domain.Lecture, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid lecture id: %d", id)
	}

	const query = `SELECT id, university, title, english_title, department, lecture_type, code, level, credit, year, language, url, abstract, goal, experience, flow, out_of_class_work, textbook, reference_book, assessment, prerequisite, contact, office_hours, note FROM lectures WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var (
		lecture                       domain.Lecture
		englishTitle, department      sql.NullString
		lectureType, language         sql.NullString
		url, abstractText             sql.NullString
		goal, experience              sql.NullString
		flow, outOfClassWork          sql.NullString
		textbook, referenceBook       sql.NullString
		assessment, prerequisite      sql.NullString
		contact, officeHours, note    sql.NullString
		levelValue, creditValue, year sql.NullInt64
	)

	err := row.Scan(
		&lecture.ID,
		&lecture.University,
		&lecture.Title,
		&englishTitle,
		&department,
		&lectureType,
		&lecture.Code,
		&levelValue,
		&creditValue,
		&year,
		&language,
		&url,
		&abstractText,
		&goal,
		&experience,
		&flow,
		&outOfClassWork,
		&textbook,
		&referenceBook,
		&assessment,
		&prerequisite,
		&contact,
		&officeHours,
		&note,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("select lecture: %w", err)
	}

	if englishTitle.Valid {
		lecture.EnglishTitle = englishTitle.String
	}
	if department.Valid {
		lecture.Department = department.String
	}
	if lectureType.Valid {
		lecture.LectureType = domain.LectureType(lectureType.String)
	}
	if levelValue.Valid {
		lecture.Level = domain.Level(levelValue.Int64)
	}
	if creditValue.Valid {
		lecture.Credit = int(creditValue.Int64)
	}
	if year.Valid {
		lecture.Year = int(year.Int64)
	}
	if language.Valid {
		lecture.Language = language.String
	}
	if url.Valid {
		lecture.Url = url.String
	}
	if abstractText.Valid {
		lecture.Abstract = abstractText.String
	}
	if goal.Valid {
		lecture.Goal = goal.String
	}
	if experience.Valid {
		lecture.Experience = experience.String
	}
	if flow.Valid {
		lecture.Flow = flow.String
	}
	if outOfClassWork.Valid {
		lecture.OutOfClassWork = outOfClassWork.String
	}
	if textbook.Valid {
		lecture.Textbook = textbook.String
	}
	if referenceBook.Valid {
		lecture.ReferenceBook = referenceBook.String
	}
	if assessment.Valid {
		lecture.Assessment = assessment.String
	}
	if prerequisite.Valid {
		lecture.Prerequisite = prerequisite.String
	}
	if contact.Valid {
		lecture.Contact = contact.String
	}
	if officeHours.Valid {
		lecture.OfficeHours = officeHours.String
	}
	if note.Valid {
		lecture.Note = note.String
	}

	timetables, err := r.fetchTimetablesMap([]int{lecture.ID})
	if err != nil {
		return nil, err
	}
	if ts, ok := timetables[lecture.ID]; ok {
		lecture.Timetables = ts
	}

	teachers, err := r.fetchTeachersMap([]int{lecture.ID})
	if err != nil {
		return nil, err
	}
	if ts, ok := teachers[lecture.ID]; ok {
		lecture.Teachers = ts
	}

	plans, err := r.fetchLecturePlans(lecture.ID)
	if err != nil {
		return nil, err
	}
	lecture.LecturePlans = plans

	keywords, err := r.fetchKeywords(lecture.ID)
	if err != nil {
		return nil, err
	}
	lecture.Keywords = keywords

	related, err := r.fetchRelatedCourses(lecture.ID)
	if err != nil {
		return nil, err
	}
	lecture.RelatedCourses = related

	return &lecture, nil
}

// Search retrieves lecture summaries filtered by the provided query fields.
func (r *LectureRepository) Search(query domain.SearchQuery) ([]domain.LectureSummary, error) {
	selectBuilder := strings.Builder{}
	selectBuilder.WriteString("SELECT DISTINCT l.id, l.university, l.title, l.department, l.code, l.level, l.year FROM lectures l")

	var joins []string
	var conditions []string
	var args []any

	if query.TeacherName != "" {
		joins = append(joins, "JOIN lecture_teachers lt ON lt.lecture_id = l.id JOIN teachers t ON t.id = lt.teacher_id")
		conditions = append(conditions, "LOWER(t.name) LIKE ?")
		args = append(args, "%"+strings.ToLower(query.TeacherName)+"%")
	}

	if len(query.Keywords) > 0 {
		joins = append(joins, "JOIN lecture_keywords lk ON lk.lecture_id = l.id")
		keywordPlaceholders := placeholders(len(query.Keywords))
		conditions = append(conditions, "lk.keyword IN ("+keywordPlaceholders+")")
		for _, keyword := range query.Keywords {
			args = append(args, keyword)
		}
	}

	if len(query.TimeTables) > 0 {
		joins = append(joins, "JOIN timetables tt ON tt.lecture_id = l.id")
		var timetableFilters []string
		for _, timetable := range query.TimeTables {
			if timetable.DayOfWeek == "" && timetable.Period == 0 {
				continue
			}
			timetableFilters = append(timetableFilters, "(tt.day_of_week = ? AND tt.period = ?)")
			args = append(args, string(timetable.DayOfWeek), int(timetable.Period))
		}
		if len(timetableFilters) > 0 {
			conditions = append(conditions, "("+strings.Join(timetableFilters, " OR ")+")")
		}
	}

	if query.Title != "" {
		conditions = append(conditions, "(LOWER(l.title) LIKE ? OR LOWER(IFNULL(l.english_title, '')) LIKE ?)")
		like := "%" + strings.ToLower(query.Title) + "%"
		args = append(args, like, like)
	}

	if len(query.Departments) > 0 {
		conditions = append(conditions, "l.department IN ("+placeholders(len(query.Departments))+")")
		for _, department := range query.Departments {
			args = append(args, department)
		}
	}

	if query.Year != 0 {
		conditions = append(conditions, "l.year = ?")
		args = append(args, query.Year)
	}

	if len(query.Levels) > 0 {
		conditions = append(conditions, "l.level IN ("+placeholders(len(query.Levels))+")")
		for _, level := range query.Levels {
			args = append(args, int(level))
		}
	}

	if len(joins) > 0 {
		selectBuilder.WriteString(" ")
		selectBuilder.WriteString(strings.Join(joins, " "))
	}

	if len(conditions) > 0 {
		selectBuilder.WriteString(" WHERE ")
		selectBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	selectBuilder.WriteString(" ORDER BY l.year DESC, l.title ASC")

	rows, err := r.db.Query(selectBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("search lectures: %w", err)
	}
	defer rows.Close()

	summaries := make([]domain.LectureSummary, 0)
	ids := make([]int, 0)

	for rows.Next() {
		var summary domain.LectureSummary
		var levelValue, yearValue sql.NullInt64
		if err := rows.Scan(&summary.ID, &summary.University, &summary.Title, &summary.Department, &summary.Code, &levelValue, &yearValue); err != nil {
			return nil, fmt.Errorf("scan lecture summary: %w", err)
		}
		if levelValue.Valid {
			summary.Level = domain.Level(levelValue.Int64)
		}
		if yearValue.Valid {
			summary.Year = int(yearValue.Int64)
		}

		summaries = append(summaries, summary)
		ids = append(ids, summary.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lecture summaries: %w", err)
	}

	if len(ids) == 0 {
		return summaries, nil
	}

	timetables, err := r.fetchTimetablesMap(ids)
	if err != nil {
		return nil, err
	}

	teachers, err := r.fetchTeachersMap(ids)
	if err != nil {
		return nil, err
	}

	for index := range summaries {
		if ts, ok := timetables[summaries[index].ID]; ok {
			summaries[index].Timetables = ts
		}
		if ts, ok := teachers[summaries[index].ID]; ok {
			summaries[index].Teachers = ts
		}
	}

	return summaries, nil
}

// Create stores a single lecture aggregate.
func (r *LectureRepository) Create(lecture *domain.Lecture) error {
	return ErrNotImplemented
}

// Creates stores multiple lecture aggregates within a single transaction.
func (r *LectureRepository) Creates(lectures []domain.Lecture) error {
	return ErrNotImplemented
}

// Update updates an existing lecture aggregate.
func (r *LectureRepository) Update(lecture *domain.Lecture) error {
	return ErrNotImplemented
}

// Delete removes a lecture and its associated records.
func (r *LectureRepository) Delete(id int) error {
	return ErrNotImplemented
}

func (r *LectureRepository) initSchema() error {
	for _, statement := range schemaStatements() {
		if _, err := r.db.Exec(statement); err != nil {
			return fmt.Errorf("init schema: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) fetchTimetablesMap(lectureIDs []int) (map[int][]domain.TimeTable, error) {
	result := make(map[int][]domain.TimeTable)
	if len(lectureIDs) == 0 {
		return result, nil
	}

	ph := placeholders(len(lectureIDs))
	query := fmt.Sprintf(`SELECT tt.lecture_id, tt.semester, tt.room_id, r.name, tt.day_of_week, tt.period FROM timetables tt LEFT JOIN rooms r ON r.id = tt.room_id WHERE tt.lecture_id IN (%s) ORDER BY tt.lecture_id, tt.semester, tt.day_of_week, tt.period`, ph)
	args := make([]any, len(lectureIDs))
	for index, id := range lectureIDs {
		args[index] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("select timetables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			lectureID           int
			semester, dayOfWeek sql.NullString
			roomID              sql.NullInt64
			roomName            sql.NullString
			period              sql.NullInt64
		)

		if err := rows.Scan(&lectureID, &semester, &roomID, &roomName, &dayOfWeek, &period); err != nil {
			return nil, fmt.Errorf("scan timetable: %w", err)
		}

		timetable := domain.TimeTable{LectureID: lectureID}
		if semester.Valid {
			timetable.Semester = domain.Semester(semester.String)
		}
		if dayOfWeek.Valid {
			timetable.DayOfWeek = domain.DayOfWeek(dayOfWeek.String)
		}
		if period.Valid {
			timetable.Period = domain.Period(period.Int64)
		}
		if roomID.Valid {
			timetable.Room.ID = int(roomID.Int64)
		}
		if roomName.Valid {
			timetable.Room.Name = roomName.String
		}

		result[lectureID] = append(result[lectureID], timetable)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate timetables: %w", err)
	}

	return result, nil
}

func (r *LectureRepository) fetchTeachersMap(lectureIDs []int) (map[int][]domain.Teacher, error) {
	result := make(map[int][]domain.Teacher)
	if len(lectureIDs) == 0 {
		return result, nil
	}

	ph := placeholders(len(lectureIDs))
	query := fmt.Sprintf(`SELECT lt.lecture_id, t.id, t.name, t.url FROM lecture_teachers lt JOIN teachers t ON t.id = lt.teacher_id WHERE lt.lecture_id IN (%s) ORDER BY lt.lecture_id, t.id`, ph)
	args := make([]any, len(lectureIDs))
	for index, id := range lectureIDs {
		args[index] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("select teachers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			lectureID int
			teacher   domain.Teacher
		)

		if err := rows.Scan(&lectureID, &teacher.ID, &teacher.Name, &teacher.Url); err != nil {
			return nil, fmt.Errorf("scan teacher: %w", err)
		}

		result[lectureID] = append(result[lectureID], teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teachers: %w", err)
	}

	return result, nil
}

func (r *LectureRepository) fetchLecturePlans(lectureID int) ([]domain.LecturePlan, error) {
	rows, err := r.db.Query(`SELECT count, plan, assignment FROM lecture_plans WHERE lecture_id = ? ORDER BY count`, lectureID)
	if err != nil {
		return nil, fmt.Errorf("select lecture plans: %w", err)
	}
	defer rows.Close()

	plans := make([]domain.LecturePlan, 0)
	for rows.Next() {
		var plan domain.LecturePlan
		if err := rows.Scan(&plan.Count, &plan.Plan, &plan.Assignment); err != nil {
			return nil, fmt.Errorf("scan lecture plan: %w", err)
		}
		plans = append(plans, plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lecture plans: %w", err)
	}

	return plans, nil
}

func (r *LectureRepository) fetchKeywords(lectureID int) ([]string, error) {
	rows, err := r.db.Query(`SELECT keyword FROM lecture_keywords WHERE lecture_id = ? ORDER BY keyword`, lectureID)
	if err != nil {
		return nil, fmt.Errorf("select keywords: %w", err)
	}
	defer rows.Close()

	keywords := make([]string, 0)
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate keywords: %w", err)
	}

	return keywords, nil
}

func (r *LectureRepository) fetchRelatedCourses(lectureID int) ([]int, error) {
	rows, err := r.db.Query(`SELECT related_lecture_id FROM related_courses WHERE lecture_id = ? ORDER BY related_lecture_id`, lectureID)
	if err != nil {
		return nil, fmt.Errorf("select related courses: %w", err)
	}
	defer rows.Close()

	related := make([]int, 0)
	for rows.Next() {
		var relatedID int
		if err := rows.Scan(&relatedID); err != nil {
			return nil, fmt.Errorf("scan related course: %w", err)
		}
		related = append(related, relatedID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate related courses: %w", err)
	}

	return related, nil
}

func placeholders(count int) string {
	if count <= 0 {
		return ""
	}

	builder := strings.Builder{}
	for i := 0; i < count; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString("?")
	}

	return builder.String()
}
