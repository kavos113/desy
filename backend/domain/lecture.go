package domain

type Lecture struct {
	ID             int
	University     string
	Title          string
	EnglishTitle   string
	Department     string
	LectureType    LectureType
	Code           string
	Credit         int
	Year           int
	Language       string
	Url            string
	Abstract       string
	Goal           string
	Experience     string
	Flow           string
	OutOfClassWork string
	Textbook       string
	ReferenceBook  string
	Assessment     string
	Prerequisite   string
	Contact        string
	OfficeHours    string
	Note           string
	Timetables     []TimeTable
	Teachers       []Teacher
	LecturePlans   []LecturePlan
	Keywords       []string
	RelatedCourses []int
}

type LectureType string

const (
	LectureTypeOffline  LectureType = "offline"
	LectureTypeLive     LectureType = "live"
	LectureTypeHyflex   LectureType = "hyflex"
	LectureTypeOndemand LectureType = "ondemand"
	LectureTypeOther    LectureType = "other"
)

type LecturePlan struct {
	Count      int
	Plan       string
	Assignment string
}
