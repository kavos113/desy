package domain

type Lecture struct {
	ID             int
	University     string
	Title          string
	EnglishTitle   string
	Department     string
	LectureType    LectureType
	Code           string
	Level          Level
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

type Level int

const (
	LevelBachelor1 Level = 1
	LevelBachelor2 Level = 2
	LevelBachelor3 Level = 3
	LevelMaster1   Level = 5
	LevelMaster2   Level = 6
	LevelDoctor    Level = 8
)

type LecturePlan struct {
	Count      int
	Plan       string
	Assignment string
}

type SearchQuery struct {
	Title       string
	Keywords    []string
	Departments []string
	Year        int
	TeacherName string
	TimeTables  []TimeTable
	Levels      []Level
}

type LectureRepository interface {
	FindByID(id int) (*Lecture, error)
	Search(query SearchQuery) ([]Lecture, error)
	Create(lecture *Lecture) error
	Creates(lectures []Lecture) error
	Update(lecture *Lecture) error	
	Delete(id int) error
}
