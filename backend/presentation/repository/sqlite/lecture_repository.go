package sqlite

import (
	"context"
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
	codes, err := r.fetchRelatedCourseCodes(lecture.ID)
	if err != nil {
		return nil, err
	}
	lecture.RelatedCourseCodes = codes

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
	if lecture == nil {
		return errors.New("nil lecture")
	}

	copies := []domain.Lecture{*lecture}
	if err := r.Creates(copies); err != nil {
		return err
	}

	lecture.ID = copies[0].ID
	lecture.Teachers = copies[0].Teachers
	lecture.Timetables = copies[0].Timetables
	lecture.LecturePlans = copies[0].LecturePlans
	lecture.Keywords = copies[0].Keywords
	lecture.RelatedCourseCodes = copies[0].RelatedCourseCodes
	lecture.RelatedCourses = copies[0].RelatedCourses

	return nil
}

// Creates stores multiple lecture aggregates within a single transaction.
func (r *LectureRepository) Creates(lectures []domain.Lecture) error {
	if len(lectures) == 0 {
		return nil
	}

	for idx := range lectures {
		lectures[idx].RelatedCourseCodes = sanitizeRelatedCourseCodes(lectures[idx].RelatedCourseCodes)
	}

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	insertedIDs := make([]int, len(lectures))
	codeToID := make(map[string]int)

	for idx := range lectures {
		id, err := r.insertLectureTx(tx, &lectures[idx])
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("rollback on insert lecture: %v (original error: %w)", rbErr, err)
			}
			return err
		}
		insertedIDs[idx] = id

		if code := normalizeCourseCode(lectures[idx].Code); code != "" {
			if _, exists := codeToID[code]; !exists {
				codeToID[code] = id
			}
		}
	}

	pendingCodes := make(map[string]struct{})
	for _, lecture := range lectures {
		for _, code := range lecture.RelatedCourseCodes {
			normalized := normalizeCourseCode(code)
			if normalized == "" {
				continue
			}
			if _, ok := codeToID[normalized]; ok {
				continue
			}
			pendingCodes[normalized] = struct{}{}
		}
	}

	if len(pendingCodes) > 0 {
		codes := make([]string, 0, len(pendingCodes))
		for code := range pendingCodes {
			codes = append(codes, code)
		}
		existing, err := r.findLectureIDsByCodesTx(tx, codes)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("rollback on resolve related codes: %v (original error: %w)", rbErr, err)
			}
			return err
		}
		for code, id := range existing {
			if _, exists := codeToID[code]; !exists {
				codeToID[code] = id
			}
		}
	}

	for idx := range lectures {
		lectureID := insertedIDs[idx]
		relatedIDs := resolveRelatedCourseIDs(lectures[idx].RelatedCourseCodes, codeToID, lectureID)
		if len(relatedIDs) > 0 {
			if err := r.insertRelatedCoursesTx(tx, lectureID, relatedIDs); err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					return fmt.Errorf("rollback on insert related courses: %v (original error: %w)", rbErr, err)
				}
				return err
			}
		}
		lectures[idx].RelatedCourses = relatedIDs
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit lecture transaction: %w", err)
	}

	for idx := range lectures {
		lectures[idx].ID = insertedIDs[idx]
	}

	return nil
}

// MigrateRelatedCourses resolves stored related course codes into lecture ID links.
func (r *LectureRepository) MigrateRelatedCourses(ctx context.Context) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin migrate related courses transaction: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT lecture_id, code FROM related_course_codes`)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, fmt.Errorf("rollback on select related course codes: %v (original error: %w)", rbErr, err)
		}
		return 0, fmt.Errorf("select related course codes: %w", err)
	}

	type pendingMapping struct {
		lectureID  int
		normalized string
	}

	mappings := make([]pendingMapping, 0)
	codeSet := make(map[string]struct{})

	for rows.Next() {
		var (
			lectureID int
			code      string
		)
		if err := rows.Scan(&lectureID, &code); err != nil {
			rows.Close()
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("rollback on scan related course code: %v (original error: %w)", rbErr, err)
			}
			return 0, fmt.Errorf("scan related course code: %w", err)
		}

		normalized := normalizeCourseCode(code)
		if normalized == "" {
			continue
		}

		mappings = append(mappings, pendingMapping{lectureID: lectureID, normalized: normalized})
		codeSet[normalized] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, fmt.Errorf("rollback after iterating related course codes: %v (original error: %w)", rbErr, err)
		}
		return 0, fmt.Errorf("iterate related course codes: %w", err)
	}
	rows.Close()

	if len(mappings) == 0 || len(codeSet) == 0 {
		if err := tx.Commit(); err != nil {
			return 0, fmt.Errorf("commit migrate related courses transaction: %w", err)
		}
		return 0, nil
	}

	codes := make([]string, 0, len(codeSet))
	for code := range codeSet {
		codes = append(codes, code)
	}

	codeToID, err := r.findLectureIDsByCodesTx(tx, codes)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, fmt.Errorf("rollback on resolve lecture ids: %v (original error: %w)", rbErr, err)
		}
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO related_courses (lecture_id, related_lecture_id) VALUES (?, ?)`)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, fmt.Errorf("rollback on prepare insert related course: %v (original error: %w)", rbErr, err)
		}
		return 0, fmt.Errorf("prepare insert related course: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	type pair struct{ src, dst int }
	seenPairs := make(map[pair]struct{})

	for _, mapping := range mappings {
		targetID, ok := codeToID[mapping.normalized]
		if !ok {
			continue
		}
		if targetID == mapping.lectureID {
			continue
		}

		key := pair{src: mapping.lectureID, dst: targetID}
		if _, exists := seenPairs[key]; exists {
			continue
		}
		seenPairs[key] = struct{}{}

		result, err := stmt.Exec(mapping.lectureID, targetID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("rollback on execute insert related course: %v (original error: %w)", rbErr, err)
			}
			return 0, fmt.Errorf("insert migrated related course: %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, fmt.Errorf("rollback on rows affected: %v (original error: %w)", rbErr, err)
			}
			return 0, fmt.Errorf("rows affected on insert related course: %w", err)
		}
		if affected > 0 {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit migrate related courses transaction: %w", err)
	}

	return inserted, nil
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
			lectureID   int
			teacherID   int
			teacherName string
			urlValue    sql.NullString
		)

		if err := rows.Scan(&lectureID, &teacherID, &teacherName, &urlValue); err != nil {
			return nil, fmt.Errorf("scan teacher: %w", err)
		}

		teacher := domain.Teacher{ID: teacherID, Name: teacherName}
		if urlValue.Valid {
			teacher.Url = urlValue.String
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

func (r *LectureRepository) fetchRelatedCourseCodes(lectureID int) ([]string, error) {
	rows, err := r.db.Query(`SELECT code FROM related_course_codes WHERE lecture_id = ? ORDER BY code`, lectureID)
	if err != nil {
		return nil, fmt.Errorf("select related course codes: %w", err)
	}
	defer rows.Close()

	codes := make([]string, 0)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("scan related course code: %w", err)
		}
		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate related course codes: %w", err)
	}

	if len(codes) == 0 {
		return nil, nil
	}

	return codes, nil
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

func (r *LectureRepository) insertLectureTx(tx *sql.Tx, lecture *domain.Lecture) (int, error) {
	if lecture == nil {
		return 0, errors.New("nil lecture")
	}
	if strings.TrimSpace(lecture.University) == "" {
		return 0, errors.New("lecture university is required")
	}
	if strings.TrimSpace(lecture.Title) == "" {
		return 0, errors.New("lecture title is required")
	}

	const insertLecture = `INSERT INTO lectures (university, title, english_title, department, lecture_type, code, level, credit, year, language, url, abstract, goal, experience, flow, out_of_class_work, textbook, reference_book, assessment, prerequisite, contact, office_hours, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(insertLecture,
		strings.TrimSpace(lecture.University),
		strings.TrimSpace(lecture.Title),
		nullString(lecture.EnglishTitle),
		nullString(lecture.Department),
		nullString(string(lecture.LectureType)),
		nullString(lecture.Code),
		nullInt(int(lecture.Level)),
		nullInt(lecture.Credit),
		nullInt(lecture.Year),
		nullString(lecture.Language),
		nullString(lecture.Url),
		nullString(lecture.Abstract),
		nullString(lecture.Goal),
		nullString(lecture.Experience),
		nullString(lecture.Flow),
		nullString(lecture.OutOfClassWork),
		nullString(lecture.Textbook),
		nullString(lecture.ReferenceBook),
		nullString(lecture.Assessment),
		nullString(lecture.Prerequisite),
		nullString(lecture.Contact),
		nullString(lecture.OfficeHours),
		nullString(lecture.Note),
	)
	if err != nil {
		return 0, fmt.Errorf("insert lecture: %w", err)
	}

	lectureID64, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}
	lectureID := int(lectureID64)

	if err := r.insertTeachersTx(tx, lectureID, lecture); err != nil {
		return 0, err
	}
	if err := r.insertTimetablesTx(tx, lectureID, lecture); err != nil {
		return 0, err
	}
	if err := r.insertLecturePlansTx(tx, lectureID, lecture.LecturePlans); err != nil {
		return 0, err
	}
	if err := r.insertKeywordsTx(tx, lectureID, lecture.Keywords); err != nil {
		return 0, err
	}
	if err := r.insertRelatedCourseCodesTx(tx, lectureID, lecture.RelatedCourseCodes); err != nil {
		return 0, err
	}

	return lectureID, nil
}

func (r *LectureRepository) insertTeachersTx(tx *sql.Tx, lectureID int, lecture *domain.Lecture) error {
	if len(lecture.Teachers) == 0 {
		return nil
	}

	seenIDs := make(map[int]struct{})
	seenNames := make(map[string]struct{})

	for idx, teacher := range lecture.Teachers {
		name := strings.TrimSpace(teacher.Name)
		if name == "" {
			continue
		}

		teacherID := teacher.ID
		if teacherID > 0 {
			if _, ok := seenIDs[teacherID]; ok {
				continue
			}
			seenIDs[teacherID] = struct{}{}
		} else {
			key := strings.ToLower(name)
			if _, ok := seenNames[key]; ok {
				continue
			}
			seenNames[key] = struct{}{}

			id, err := r.ensureTeacherTx(tx, name, teacher.Url)
			if err != nil {
				return err
			}
			teacherID = id
			lecture.Teachers[idx].ID = id
		}

		if _, err := tx.Exec(`INSERT OR IGNORE INTO lecture_teachers (lecture_id, teacher_id) VALUES (?, ?)`, lectureID, teacherID); err != nil {
			return fmt.Errorf("insert lecture teacher: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) insertTimetablesTx(tx *sql.Tx, lectureID int, lecture *domain.Lecture) error {
	if len(lecture.Timetables) == 0 {
		return nil
	}

	for idx, timetable := range lecture.Timetables {
		roomID := 0
		if name := strings.TrimSpace(timetable.Room.Name); name != "" {
			id, err := r.ensureRoomTx(tx, name)
			if err != nil {
				return err
			}
			roomID = id
			lecture.Timetables[idx].Room.ID = id
		}

		if _, err := tx.Exec(`INSERT INTO timetables (lecture_id, semester, room_id, day_of_week, period) VALUES (?, ?, ?, ?, ?)`,
			lectureID,
			nullString(string(timetable.Semester)),
			nullInt(roomID),
			nullString(string(timetable.DayOfWeek)),
			nullInt(int(timetable.Period)),
		); err != nil {
			return fmt.Errorf("insert timetable: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) insertLecturePlansTx(tx *sql.Tx, lectureID int, plans []domain.LecturePlan) error {
	if len(plans) == 0 {
		return nil
	}

	for _, plan := range plans {
		if _, err := tx.Exec(`INSERT INTO lecture_plans (lecture_id, count, plan, assignment) VALUES (?, ?, ?, ?)`,
			lectureID,
			nullInt(plan.Count),
			nullString(plan.Plan),
			nullString(plan.Assignment),
		); err != nil {
			return fmt.Errorf("insert lecture plan: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) insertRelatedCourseCodesTx(tx *sql.Tx, lectureID int, codes []string) error {
	if _, err := tx.Exec(`DELETE FROM related_course_codes WHERE lecture_id = ?`, lectureID); err != nil {
		return fmt.Errorf("delete related course codes: %w", err)
	}

	if len(codes) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(`INSERT INTO related_course_codes (lecture_id, code) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert related course code: %w", err)
	}
	defer stmt.Close()

	for _, code := range codes {
		if _, err := stmt.Exec(lectureID, code); err != nil {
			return fmt.Errorf("insert related course code: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) insertKeywordsTx(tx *sql.Tx, lectureID int, keywords []string) error {
	if len(keywords) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	for _, keyword := range keywords {
		clean := strings.TrimSpace(keyword)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if _, err := tx.Exec(`INSERT OR IGNORE INTO lecture_keywords (lecture_id, keyword) VALUES (?, ?)`, lectureID, clean); err != nil {
			return fmt.Errorf("insert keyword: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) insertRelatedCoursesTx(tx *sql.Tx, lectureID int, related []int) error {
	if len(related) == 0 {
		return nil
	}

	seen := make(map[int]struct{})
	for _, rel := range related {
		if rel <= 0 {
			continue
		}
		if _, ok := seen[rel]; ok {
			continue
		}
		seen[rel] = struct{}{}
		if _, err := tx.Exec(`INSERT OR IGNORE INTO related_courses (lecture_id, related_lecture_id) VALUES (?, ?)`, lectureID, rel); err != nil {
			return fmt.Errorf("insert related course: %w", err)
		}
	}

	return nil
}

func (r *LectureRepository) findLectureIDsByCodesTx(tx *sql.Tx, codes []string) (map[string]int, error) {
	result := make(map[string]int)
	if len(codes) == 0 {
		return result, nil
	}

	placeholders := placeholders(len(codes))
	query := fmt.Sprintf(`SELECT id, code FROM lectures WHERE UPPER(code) IN (%s)`, placeholders)
	args := make([]any, len(codes))
	for idx, code := range codes {
		args[idx] = code
	}

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("select lectures by code: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			return nil, fmt.Errorf("scan lecture code: %w", err)
		}
		normalized := normalizeCourseCode(code)
		if normalized == "" {
			continue
		}
		if _, exists := result[normalized]; !exists {
			result[normalized] = id
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lecture codes: %w", err)
	}

	return result, nil
}

func resolveRelatedCourseIDs(codes []string, mapping map[string]int, selfID int) []int {
	if len(codes) == 0 {
		return nil
	}
	results := make([]int, 0, len(codes))
	seen := make(map[int]struct{})
	for _, code := range codes {
		normalized := normalizeCourseCode(code)
		if normalized == "" {
			continue
		}
		id, ok := mapping[normalized]
		if !ok {
			continue
		}
		if id == selfID {
			continue
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		results = append(results, id)
	}
	if len(results) == 0 {
		return nil
	}
	return results
}

func normalizeCourseCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func sanitizeRelatedCourseCodes(codes []string) []string {
	if len(codes) == 0 {
		return nil
	}

	result := make([]string, 0, len(codes))
	seen := make(map[string]struct{})
	for _, code := range codes {
		trimmed := strings.TrimSpace(code)
		if trimmed == "" {
			continue
		}
		normalized := normalizeCourseCode(trimmed)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, trimmed)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func (r *LectureRepository) ensureTeacherTx(tx *sql.Tx, name, url string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM teachers WHERE name = ?`, name).Scan(&id)
	if err == nil {
		if strings.TrimSpace(url) != "" {
			if _, err := tx.Exec(`UPDATE teachers SET url = ? WHERE id = ?`, url, id); err != nil {
				return 0, fmt.Errorf("update teacher url: %w", err)
			}
		}
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("select teacher: %w", err)
	}

	result, err := tx.Exec(`INSERT INTO teachers (name, url) VALUES (?, ?)`, name, nullString(url))
	if err != nil {
		return 0, fmt.Errorf("insert teacher: %w", err)
	}

	teacherID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert teacher id: %w", err)
	}

	return int(teacherID), nil
}

func (r *LectureRepository) ensureRoomTx(tx *sql.Tx, name string) (int, error) {
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

func nullString(value string) sql.NullString {
	if strings.TrimSpace(value) == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullInt(value int) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(value), Valid: true}
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
