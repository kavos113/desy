package scraper

import (
	"os"
	"strings"
	"testing"

	"github.com/kavos113/desy/backend/domain"
)

func TestParseCourseList(t *testing.T) {
	file, err := os.Open("fixture/course_list.html")
	if err != nil {
		t.Fatalf("open course list fixture: %v", err)
	}
	defer file.Close()

	items, err := ParseCourseList(file, "https://syllabus.s.isct.ac.jp")
	if err != nil {
		t.Fatalf("ParseCourseList returned error: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	first := items[0]
	if first.Code != "MCS.M301" {
		t.Errorf("unexpected code: %s", first.Code)
	}
	if first.Title != "B2Dプレ研究実践a（数理・計算科学系）" {
		t.Errorf("unexpected title: %s", first.Title)
	}
	if first.DetailURL != "https://syllabus.s.isct.ac.jp/courses/2025/4/0-904-342200-0-0/202534819" {
		t.Errorf("unexpected url: %s", first.DetailURL)
	}

	third := items[2]
	if third.Code != "MCS.T201" {
		t.Errorf("unexpected third code: %s", third.Code)
	}
	if third.DetailURL == "" {
		t.Fatalf("third detail url should not be empty")
	}
}

func TestParseCourseDetail(t *testing.T) {
	file, err := os.Open("fixture/course_detail.html")
	if err != nil {
		t.Fatalf("open course detail fixture: %v", err)
	}
	defer file.Close()

	lecture, err := ParseCourseDetail(file, "https://example.com/courses/2025/LAH.S101")
	if err != nil {
		t.Fatalf("ParseCourseDetail returned error: %v", err)
	}

	expected := &domain.Lecture{
		ID: 	   1,
		University: "東京科学大学",
		Title:      "法学（憲法）Ａ",
		Department: "文系教養科目",
		LectureType: domain.LectureTypeOffline,	
		Code:       "LAH.S101",
		Level: domain.LevelBachelor1,
		Credit: 1,
		Year:       2025,
		Language: "日本語",
		Url: "https://example.com/courses/2025/LAH.S101",
		Abstract: "憲法の基本的人権",
		Goal: "基本的人権",
		Experience: "弁護士",
		Flow: "法律初学者を想定し",
		OutOfClassWork: "学修効果",
		Textbook: "憲法（第八版）",
		ReferenceBook: "マンガ僕たちの日本国憲法",
		Assessment: "期末",
		Prerequisite: "歓迎します",
		Contact: "",
		OfficeHours: "",
		Note: "日本国憲法は",
		Timetables: []domain.TimeTable{
			{
				LectureID: 1,
				Semester: domain.SemesterFall,
				DayOfWeek: domain.DayOfWeekMonday,
				Period:    domain.Period5,
				Room: domain.Room{
					Name: "S3-215(S321)",
				},
			},
			{
				LectureID: 1,
				Semester: domain.SemesterFall,
				DayOfWeek: domain.DayOfWeekMonday,
				Period:    domain.Period6,
				Room: domain.Room{
					Name: "S3-215(S321)",
				},
			},
		},
		Teachers: []domain.Teacher{
			{
				Name: "篠島 正幸",
			},
		},
		LecturePlans: []domain.LecturePlan{
			{
				Count: 1,
				Plan: "講義ガイダンス",
				Assignment: "社会における",
			},
			{
				Count: 2,
				Plan: "憲法の基本理念",
				Assignment: "憲法の条項",
			},
		},
		Keywords: []string{"憲法", "法律", "人権", "教養"},
		RelatedCourses: []int{},
	}

	if lecture.Title != expected.Title {
		t.Errorf("unexpected title: got %s, want %s", lecture.Title, expected.Title)
	}
	if lecture.University != expected.University {
		t.Errorf("unexpected university: got %s, want %s", lecture.University, expected.University)
	}
	if lecture.Title != expected.Title {
		t.Errorf("unexpected title: got %s, want %s", lecture.Title, expected.Title)
	}
	if lecture.Department != expected.Department {
		t.Errorf("unexpected department: got %s, want %s", lecture.Department, expected.Department)
	}
	if lecture.LectureType != expected.LectureType {
		t.Errorf("unexpected lecture type: got %s, want %s", lecture.LectureType, expected.LectureType)
	}
	if lecture.Code != expected.Code {
		t.Errorf("unexpected code: got %s, want %s", lecture.Code, expected.Code)
	}
	if lecture.Level != expected.Level {
		t.Errorf("unexpected level: got %d, want %d", lecture.Level, expected.Level)
	}
	if lecture.Credit != expected.Credit {
		t.Errorf("unexpected credit: got %d, want %d", lecture.Credit, expected.Credit)
	}
	if lecture.Year != expected.Year {
		t.Errorf("unexpected year: got %d, want %d", lecture.Year, expected.Year)
	}
	if lecture.Language != expected.Language {
		t.Errorf("unexpected language: got %s, want %s", lecture.Language, expected.Language)
	}
	if lecture.Url != expected.Url {
		t.Errorf("unexpected url: got %s, want %s", lecture.Url, expected.Url)
	}
	if !strings.Contains(lecture.Abstract, expected.Abstract) {
		t.Errorf("unexpected abstract: got %s, want to contain %s", lecture.Abstract, expected.Abstract)
	}
	if !strings.Contains(lecture.Goal, expected.Goal) {
		t.Errorf("unexpected goal: got %s, want to contain %s", lecture.Goal, expected.Goal)
	}
	if !strings.Contains(lecture.Experience, expected.Experience) {
		t.Errorf("unexpected experience: got %s, want to contain %s", lecture.Experience, expected.Experience)
	}
	if !strings.Contains(lecture.Flow, expected.Flow) {
		t.Errorf("unexpected flow: got %s, want to contain %s", lecture.Flow, expected.Flow)
	}
	if !strings.Contains(lecture.OutOfClassWork, expected.OutOfClassWork) {
		t.Errorf("unexpected out of class work: got %s, want to contain %s", lecture.OutOfClassWork, expected.OutOfClassWork)
	}
	if !strings.Contains(lecture.Textbook, expected.Textbook) {
		t.Errorf("unexpected textbook: got %s, want to contain %s", lecture.Textbook, expected.Textbook)
	}
	if !strings.Contains(lecture.ReferenceBook, expected.ReferenceBook) {
		t.Errorf("unexpected reference book: got %s, want to contain %s", lecture.ReferenceBook, expected.ReferenceBook)
	}
	if !strings.Contains(lecture.Assessment, expected.Assessment) {
		t.Errorf("unexpected assessment: got %s, want to contain %s", lecture.Assessment, expected.Assessment)
	}
	if !strings.Contains(lecture.Prerequisite, expected.Prerequisite) {
		t.Errorf("unexpected prerequisite: got %s, want to contain %s", lecture.Prerequisite, expected.Prerequisite)
	}
	if !strings.Contains(lecture.Contact, expected.Contact) {
		t.Errorf("unexpected contact: got %s, want to contain %s", lecture.Contact, expected.Contact)	
	}
	if !strings.Contains(lecture.OfficeHours, expected.OfficeHours) {
		t.Errorf("unexpected office hours: got %s, want to contain %s", lecture.OfficeHours, expected.OfficeHours)
	}
	if !strings.Contains(lecture.Note, expected.Note) {
		t.Errorf("unexpected note: got %s, want to contain %s", lecture.Note, expected.Note)
	}

	if len(lecture.Timetables) != len(expected.Timetables) {
		t.Fatalf("unexpected number of timetables: got %d, want %d", len(lecture.Timetables), len(expected.Timetables))
	}
	for i, tt := range lecture.Timetables {
		expTT := expected.Timetables[i]
		if tt.Semester != expTT.Semester {
			t.Errorf("unexpected timetable semester at index %d: got %s, want %s", i, tt.Semester, expTT.Semester)
		}
		if tt.DayOfWeek != expTT.DayOfWeek {
			t.Errorf("unexpected timetable day of week at index %d: got %s, want %s", i, tt.DayOfWeek, expTT.DayOfWeek)
		}
		if tt.Period != expTT.Period {
			t.Errorf("unexpected timetable period at index %d: got %d, want %d", i, tt.Period, expTT.Period)
		}
		if tt.Room.Name != expTT.Room.Name {
			t.Errorf("unexpected timetable room name at index %d: got %s, want %s", i, tt.Room.Name, expTT.Room.Name)
		}
	}

	if len(lecture.Teachers) != len(expected.Teachers) {
		t.Fatalf("unexpected number of teachers: got %d, want %d", len(lecture.Teachers), len(expected.Teachers))
	}
	for i, teacher := range lecture.Teachers {
		expTeacher := expected.Teachers[i]
		if teacher.Name != expTeacher.Name {
			t.Errorf("unexpected teacher name at index %d: got %s, want %s", i, teacher.Name, expTeacher.Name)
		}
	}

	if len(lecture.LecturePlans) != len(expected.LecturePlans) {
		t.Fatalf("unexpected number of lecture plans: got %d, want %d", len(lecture.LecturePlans), len(expected.LecturePlans))
	}
	for i, plan := range lecture.LecturePlans {
		expPlan := expected.LecturePlans[i]
		if plan.Count != expPlan.Count {
			t.Errorf("unexpected lecture plan count at index %d: got %d, want %d", i, plan.Count, expPlan.Count)
		}
		if !strings.Contains(plan.Plan, expPlan.Plan) {
			t.Errorf("unexpected lecture plan at index %d: got %s, want to contain %s", i, plan.Plan, expPlan.Plan)
		}
		if !strings.Contains(plan.Assignment, expPlan.Assignment) {
			t.Errorf("unexpected lecture assignment at index %d: got %s, want to contain %s", i, plan.Assignment, expPlan.Assignment)
		}
	}

	if len(lecture.Keywords) != len(expected.Keywords) {
		t.Fatalf("unexpected number of keywords: got %d, want %d", len(lecture.Keywords), len(expected.Keywords))
	}
	for i, keyword := range lecture.Keywords {
		expKeyword := expected.Keywords[i]
		if keyword != expKeyword {
			t.Errorf("unexpected keyword at index %d: got %s, want %s", i, keyword, expKeyword)
		}
	}

	if len(lecture.RelatedCourses) != len(expected.RelatedCourses) {
		t.Fatalf("unexpected number of related courses: got %d, want %d", len(lecture.RelatedCourses), len(expected.RelatedCourses))
	}
	for i, relatedCourse := range lecture.RelatedCourses {
		expRelatedCourse := expected.RelatedCourses[i]
		if relatedCourse != expRelatedCourse {
			t.Errorf("unexpected related course at index %d: got %d, want %d", i, relatedCourse, expRelatedCourse)
		}
	}	
}
