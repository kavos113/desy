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

	if lecture.Title != "法学（憲法）Ａ" {
		t.Errorf("unexpected title: %s", lecture.Title)
	}
	if lecture.Department != "文系教養科目" {
		t.Errorf("unexpected department: %s", lecture.Department)
	}
	if lecture.Code != "LAH.S101" {
		t.Errorf("unexpected code: %s", lecture.Code)
	}
	if lecture.Language != "日本語" {
		t.Errorf("unexpected language: %s", lecture.Language)
	}
	if lecture.Year != 2025 {
		t.Errorf("unexpected year: %d", lecture.Year)
	}
	if lecture.Credit != 1 {
		t.Errorf("unexpected credit: %d", lecture.Credit)
	}
	if lecture.LectureType != domain.LectureTypeOffline {
		t.Errorf("unexpected lecture type: %s", lecture.LectureType)
	}

	if got, want := len(lecture.Teachers), 1; got != want {
		t.Fatalf("unexpected teachers length: %d", got)
	}
	if lecture.Teachers[0].Name != "篠島 正幸" {
		t.Errorf("unexpected teacher: %s", lecture.Teachers[0].Name)
	}

	if len(lecture.Keywords) != 4 {
		t.Fatalf("unexpected keywords length: %d", len(lecture.Keywords))
	}
	if lecture.Keywords[0] != "憲法" {
		t.Errorf("unexpected first keyword: %s", lecture.Keywords[0])
	}

	if len(lecture.Timetables) != 2 {
		t.Fatalf("unexpected timetable length: %d", len(lecture.Timetables))
	}
	if lecture.Timetables[0].DayOfWeek != domain.DayOfWeekMonday {
		t.Errorf("unexpected day of week: %s", lecture.Timetables[0].DayOfWeek)
	}
	if lecture.Timetables[0].Period != domain.Period5 {
		t.Errorf("unexpected first period: %d", lecture.Timetables[0].Period)
	}
	if lecture.Timetables[0].Room.Name != "S3-215(S321)" {
		t.Errorf("unexpected room name: %s", lecture.Timetables[0].Room.Name)
	}

	if len(lecture.LecturePlans) != 8 {
		t.Fatalf("unexpected lecture plans length: %d", len(lecture.LecturePlans))
	}
	if lecture.LecturePlans[0].Count != 1 {
		t.Errorf("unexpected first plan count: %d", lecture.LecturePlans[0].Count)
	}
	if !strings.Contains(lecture.LecturePlans[0].Plan, "講義ガイダンス") {
		t.Errorf("unexpected plan text: %s", lecture.LecturePlans[0].Plan)
	}
	if !strings.Contains(lecture.LecturePlans[0].Assignment, "CaseStudy") {
		t.Errorf("unexpected assignment text: %s", lecture.LecturePlans[0].Assignment)
	}

	if !strings.Contains(lecture.Abstract, "憲法の基本的人権") {
		t.Errorf("unexpected abstract: %s", lecture.Abstract)
	}
	if !strings.Contains(lecture.Goal, "基本的人権") {
		t.Errorf("unexpected goal: %s", lecture.Goal)
	}
	if !strings.Contains(lecture.Flow, "法律初学者を想定し") {
		t.Errorf("unexpected flow: %s", lecture.Flow)
	}
	if !strings.Contains(lecture.OutOfClassWork, "学修効果") {
		t.Errorf("unexpected out of class work: %s", lecture.OutOfClassWork)
	}
	if !strings.Contains(lecture.Textbook, "憲法（第八版）") {
		t.Errorf("unexpected textbook: %s", lecture.Textbook)
	}
	if !strings.Contains(lecture.ReferenceBook, "マンガ僕たちの日本国憲法") {
		t.Errorf("unexpected reference book: %s", lecture.ReferenceBook)
	}
	if !strings.Contains(lecture.Assessment, "期末") {
		t.Errorf("unexpected assessment: %s", lecture.Assessment)
	}
	if !strings.Contains(lecture.Experience, "弁護士") {
		t.Errorf("unexpected experience text: %s", lecture.Experience)
	}
	if !strings.Contains(lecture.Prerequisite, "歓迎します") {
		t.Errorf("unexpected prerequisite: %s", lecture.Prerequisite)
	}
	if !strings.Contains(lecture.Note, "日本国憲法は") {
		t.Errorf("unexpected note: %s", lecture.Note)
	}
}
