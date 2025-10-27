package domain

import "strings"

type Semester string

const (
	SemesterSpring Semester = "spring"
	SemesterSummer Semester = "summer"
	SemesterFall   Semester = "fall"
	SemesterWinter Semester = "winter"
)

func fromQuarter(quarter string) []Semester {
	quarter = strings.TrimSpace(quarter)
	if quarter == "" {
		return nil
	}

	quarter = strings.ReplaceAll(quarter, " ", "")

	if quarter == "通年" {
		return []Semester{SemesterSpring, SemesterSummer, SemesterFall, SemesterWinter}
	}

	separators := []string{"-", "～", "〜"}
	for _, sep := range separators {
		if strings.Contains(quarter, sep) {
			parts := strings.SplitN(quarter, sep, 2)
			if len(parts) != 2 {
				return nil
			}
			for i := range parts {
				if !strings.HasSuffix(parts[i], "Q") && strings.ContainsAny(parts[i], "1234") {
					parts[i] = parts[i] + "Q"
				}
			}
			startIdx, ok := quarterIndex(parts[0])
			if !ok {
				return nil
			}
			endIdx, ok := quarterIndex(parts[1])
			if !ok || endIdx < startIdx {
				return nil
			}

			all := []Semester{SemesterSpring, SemesterSummer, SemesterFall, SemesterWinter}
			return append([]Semester{}, all[startIdx:endIdx+1]...)
		}
	}

	if idx, ok := quarterIndex(quarter); ok {
		order := []Semester{SemesterSpring, SemesterSummer, SemesterFall, SemesterWinter}
		return []Semester{order[idx]}
	}

	if !strings.HasSuffix(quarter, "Q") && strings.ContainsAny(quarter, "1234") {
		if idx, ok := quarterIndex(quarter + "Q"); ok {
			order := []Semester{SemesterSpring, SemesterSummer, SemesterFall, SemesterWinter}
			return []Semester{order[idx]}
		}
	}

	return nil
}

func quarterIndex(token string) (int, bool) {
	token = strings.TrimSpace(token)
	switch token {
	case "1Q", "春学期", "Spring":
		return 0, true
	case "2Q", "夏学期", "Summer":
		return 1, true
	case "3Q", "秋学期", "Autumn", "Fall":
		return 2, true
	case "4Q", "冬学期", "Winter":
		return 3, true
	default:
		return 0, false
	}
}

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
	Period1  Period = 1
	Period2  Period = 2
	Period3  Period = 3
	Period4  Period = 4
	Period5  Period = 5
	Period6  Period = 6
	Period7  Period = 7
	Period8  Period = 8
	Period9  Period = 9
	Period10 Period = 10
	Period11 Period = 11
	Period12 Period = 12
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
