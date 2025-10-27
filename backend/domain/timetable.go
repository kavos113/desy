package domain

type Semester string

const (
	SemesterSpring Semester = "spring"
	SemesterSummer Semester = "summer"
	SemesterFall   Semester = "fall"
	SemesterWinter Semester = "winter"
)

type DayOfWeek string

const (
	DayOfWeekMonday    DayOfWeek = "monday"
	DayOfWeekTuesday   DayOfWeek = "tuesday"
	DayOfWeekWednesday DayOfWeek = "wednesday"
	DayOfWeekThursday  DayOfWeek = "thursday"
	DayOfWeekFriday    DayOfWeek = "friday"
	DayOfWeekSaturday  DayOfWeek = "saturday"
	DayOfWeekSunday    DayOfWeek = "sunday"
)

type Period int

const (
	Period1 Period = 1
	Period2 Period = 2
	Period3 Period = 3
	Period4 Period = 4
	Period5 Period = 5
	Period6 Period = 6
	Period7 Period = 7
	Period8 Period = 8
)

type TimeTable struct {
	LectureID int
	Semester  Semester
	Room      Room
	DayOfWeek DayOfWeek
	Period    Period
}

type Room struct {
	ID   int
	Name string
}

type TimeTableRepository interface {
	FindByLectureID(lectureID int) ([]TimeTable, error)
	Create(timetable *TimeTable) error
	Creates(timetables []TimeTable) error
	Update(timetable *TimeTable) error
	Delete(lectureID int) error
}