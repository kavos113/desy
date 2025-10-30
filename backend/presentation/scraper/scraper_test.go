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

	f, err := os.Open("fixture/course_list_full.html")
	if err != nil {
		t.Fatalf("open course list fixture: %v", err)
	}
	defer f.Close()

	items, err = ParseCourseList(f, "https://syllabus.s.isct.ac.jp")
	if err != nil {
		t.Fatalf("ParseCourseList returned error: %v", err)
	}

	first = items[0]
	if first.Code != "MCS.M201" {
		t.Errorf("unexpected code: %s", first.Code)
	}
	if first.Title != "B2D研究基礎a（数理・計算科学系）" {
		t.Errorf("unexpected title: %s", first.Title)
	}
	if first.DetailURL != "https://syllabus.s.isct.ac.jp/courses/2025/4/0-904-342200-0-0/202534815" {
		t.Errorf("unexpected url: %s", first.DetailURL)
	}
}

func TestParseCourseDetail(t *testing.T) {
	testCases := []struct {
		name      string
		fixture   string
		detailURL string
		expected  domain.Lecture
	}{
		{
			name:      "default",
			fixture:   "fixture/course_detail.html",
			detailURL: "https://example.com/courses/2025/LAH.S101",
			expected: domain.Lecture{
				ID:             1,
				University:     "東京科学大学",
				Title:          "法学（憲法）Ａ",
				Department:     "文系教養科目",
				LectureType:    domain.LectureTypeOffline,
				Code:           "LAH.S101",
				Level:          domain.LevelBachelor1,
				Credit:         1,
				Year:           2025,
				Language:       "日本語",
				Url:            "https://example.com/courses/2025/LAH.S101",
				Abstract:       "憲法の基本的人権",
				Goal:           "基本的人権",
				Experience:     "弁護士",
				Flow:           "法律初学者を想定し",
				OutOfClassWork: "学修効果",
				Textbook:       "憲法（第八版）",
				ReferenceBook:  "マンガ僕たちの日本国憲法",
				Assessment:     "期末",
				Prerequisite:   "歓迎します",
				Contact:        "",
				OfficeHours:    "",
				Note:           "日本国憲法は",
				Timetables: []domain.TimeTable{
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
						DayOfWeek: domain.DayOfWeekMonday,
						Period:    domain.Period5,
						Room: domain.Room{
							Name: "S3-215(S321)",
						},
					},
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
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
						Count:      1,
						Plan:       "講義ガイダンス",
						Assignment: "社会における",
					},
					{
						Count:      2,
						Plan:       "憲法の基本理念",
						Assignment: "憲法の条項",
					},
				},
				Keywords:           []string{"憲法", "法律", "人権", "教養"},
				RelatedCourseCodes: []string{"LAH.S201", "LAH.S301", "LAH.S102"},
				RelatedCourses:     []int{},
			},
		},
		{
			name:      "variant",
			fixture:   "fixture/course_detail2.html",
			detailURL: "https://example.com/courses/2025/LAH.S101",
			expected: domain.Lecture{
				ID:             1,
				University:     "東京科学大学",
				Title:          "法学（憲法）Ａ",
				Department:     "文系教養科目",
				LectureType:    domain.LectureTypeOffline,
				Code:           "LAH.S101",
				Level:          domain.LevelBachelor1,
				Credit:         1,
				Year:           2025,
				Language:       "日本語",
				Url:            "https://example.com/courses/2025/LAH.S101",
				Abstract:       "憲法の基本的人権",
				Goal:           "基本的人権",
				Experience:     "弁護士",
				Flow:           "法律初学者を想定し",
				OutOfClassWork: "学修効果",
				Textbook:       "憲法（第八版）",
				ReferenceBook:  "マンガ僕たちの日本国憲法",
				Assessment:     "期末",
				Prerequisite:   "",
				Contact:        "金子晴彦: kaneko[at]c.titech.ac.jp",
				OfficeHours:    "メールで事前予約すること。",
				Note:           "",
				Timetables: []domain.TimeTable{
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
						DayOfWeek: domain.DayOfWeekMonday,
						Period:    domain.Period5,
						Room: domain.Room{
							Name: "S3-215(S321)",
						},
					},
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
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
						Count:      1,
						Plan:       "講義ガイダンス",
						Assignment: "社会における",
					},
					{
						Count:      2,
						Plan:       "憲法の基本理念",
						Assignment: "憲法の条項",
					},
				},
				Keywords:           []string{"憲法", "法律", "人権", "教養"},
				RelatedCourseCodes: []string{"LAH.S201", "LAH.S301", "LAH.S102"},
				RelatedCourses:     []int{},
			},
		},
		{
			name:      "default",
			fixture:   "fixture/course_detail_full.html",
			detailURL: "https://example.com/courses/2025/LAH.S101",
			expected: domain.Lecture{
				ID:             1,
				University:     "東京科学大学",
				Title:          "法学（憲法）Ａ",
				Department:     "文系教養科目",
				LectureType:    domain.LectureTypeOffline,
				Code:           "LAH.S101",
				Level:          domain.LevelBachelor1,
				Credit:         1,
				Year:           2025,
				Language:       "日本語",
				Url:            "https://example.com/courses/2025/LAH.S101",
				Abstract:       "憲法の基本的人権",
				Goal:           "基本的人権",
				Experience:     "弁護士",
				Flow:           "法律初学者を想定し",
				OutOfClassWork: "学修効果",
				Textbook:       "憲法（第八版）",
				ReferenceBook:  "マンガ僕たちの日本国憲法",
				Assessment:     "期末",
				Prerequisite:   "歓迎します",
				Contact:        "",
				OfficeHours:    "",
				Note:           "日本国憲法は",
				Timetables: []domain.TimeTable{
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
						DayOfWeek: domain.DayOfWeekMonday,
						Period:    domain.Period5,
						Room: domain.Room{
							Name: "S3-215(S321)",
						},
					},
					{
						LectureID: 1,
						Semester:  domain.SemesterFall,
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
						Count:      1,
						Plan:       "講義ガイダンス",
						Assignment: "社会における",
					},
					{
						Count:      2,
						Plan:       "憲法の基本理念",
						Assignment: "憲法の条項",
					},
					{
						Count:      3,
						Plan:       "人権総論",
						Assignment: "人権に関する",
					},
					{
						Count:      4,
						Plan:       "法の下の平等",
						Assignment: "法の下の平等",
					},
					{
						Count:      5,
						Plan:       "精神的自由権",
						Assignment: "表現の自由",
					},
					{
						Count:      6,
						Plan:       "経済的自由権",
						Assignment: "社会経済",
					},
					{
						Count:      7,
						Plan:       "国務請求権",
						Assignment: "各人権",
					},
					{
						Count:      8,
						Plan:       "総括と試験",
						Assignment: "",
					},
				},
				Keywords:           []string{"憲法", "法律", "人権", "教養"},
				RelatedCourseCodes: []string{"LAH.S201", "LAH.S301", "LAH.S102"},
				RelatedCourses:     []int{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.fixture)
			if err != nil {
				t.Fatalf("open course detail fixture: %v", err)
			}
			defer file.Close()

			lecture, err := ParseCourseDetail(file, tc.detailURL)
			if err != nil {
				t.Fatalf("ParseCourseDetail returned error: %v", err)
			}

			if lecture.Title != tc.expected.Title {
				t.Errorf("unexpected title: got %s, want %s", lecture.Title, tc.expected.Title)
			}
			if lecture.University != tc.expected.University {
				t.Errorf("unexpected university: got %s, want %s", lecture.University, tc.expected.University)
			}
			if lecture.Department != tc.expected.Department {
				t.Errorf("unexpected department: got %s, want %s", lecture.Department, tc.expected.Department)
			}
			if lecture.LectureType != tc.expected.LectureType {
				t.Errorf("unexpected lecture type: got %s, want %s", lecture.LectureType, tc.expected.LectureType)
			}
			if lecture.Code != tc.expected.Code {
				t.Errorf("unexpected code: got %s, want %s", lecture.Code, tc.expected.Code)
			}
			if lecture.Level != tc.expected.Level {
				t.Errorf("unexpected level: got %d, want %d", lecture.Level, tc.expected.Level)
			}
			if lecture.Credit != tc.expected.Credit {
				t.Errorf("unexpected credit: got %d, want %d", lecture.Credit, tc.expected.Credit)
			}
			if lecture.Year != tc.expected.Year {
				t.Errorf("unexpected year: got %d, want %d", lecture.Year, tc.expected.Year)
			}
			if lecture.Language != tc.expected.Language {
				t.Errorf("unexpected language: got %s, want %s", lecture.Language, tc.expected.Language)
			}
			if lecture.Url != tc.expected.Url {
				t.Errorf("unexpected url: got %s, want %s", lecture.Url, tc.expected.Url)
			}
			if !strings.Contains(lecture.Abstract, tc.expected.Abstract) {
				t.Errorf("unexpected abstract: got %s, want to contain %s", lecture.Abstract, tc.expected.Abstract)
			}
			if !strings.Contains(lecture.Goal, tc.expected.Goal) {
				t.Errorf("unexpected goal: got %s, want to contain %s", lecture.Goal, tc.expected.Goal)
			}
			if !strings.Contains(lecture.Experience, tc.expected.Experience) {
				t.Errorf("unexpected experience: got %s, want to contain %s", lecture.Experience, tc.expected.Experience)
			}
			if !strings.Contains(lecture.Flow, tc.expected.Flow) {
				t.Errorf("unexpected flow: got %s, want to contain %s", lecture.Flow, tc.expected.Flow)
			}
			if !strings.Contains(lecture.OutOfClassWork, tc.expected.OutOfClassWork) {
				t.Errorf("unexpected out of class work: got %s, want to contain %s", lecture.OutOfClassWork, tc.expected.OutOfClassWork)
			}
			if !strings.Contains(lecture.Textbook, tc.expected.Textbook) {
				t.Errorf("unexpected textbook: got %s, want to contain %s", lecture.Textbook, tc.expected.Textbook)
			}
			if !strings.Contains(lecture.ReferenceBook, tc.expected.ReferenceBook) {
				t.Errorf("unexpected reference book: got %s, want to contain %s", lecture.ReferenceBook, tc.expected.ReferenceBook)
			}
			if !strings.Contains(lecture.Assessment, tc.expected.Assessment) {
				t.Errorf("unexpected assessment: got %s, want to contain %s", lecture.Assessment, tc.expected.Assessment)
			}
			if !strings.Contains(lecture.Prerequisite, tc.expected.Prerequisite) {
				t.Errorf("unexpected prerequisite: got %s, want to contain %s", lecture.Prerequisite, tc.expected.Prerequisite)
			}
			if !strings.Contains(lecture.Contact, tc.expected.Contact) {
				t.Errorf("unexpected contact: got %s, want to contain %s", lecture.Contact, tc.expected.Contact)
			}
			if !strings.Contains(lecture.OfficeHours, tc.expected.OfficeHours) {
				t.Errorf("unexpected office hours: got %s, want to contain %s", lecture.OfficeHours, tc.expected.OfficeHours)
			}
			if !strings.Contains(lecture.Note, tc.expected.Note) {
				t.Errorf("unexpected note: got %s, want to contain %s", lecture.Note, tc.expected.Note)
			}

			if len(lecture.Timetables) != len(tc.expected.Timetables) {
				t.Fatalf("unexpected number of timetables: got %d, want %d", len(lecture.Timetables), len(tc.expected.Timetables))
			}
			for i, tt := range lecture.Timetables {
				expTT := tc.expected.Timetables[i]
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

			if len(lecture.Teachers) != len(tc.expected.Teachers) {
				t.Fatalf("unexpected number of teachers: got %d, want %d", len(lecture.Teachers), len(tc.expected.Teachers))
			}
			for i, teacher := range lecture.Teachers {
				expTeacher := tc.expected.Teachers[i]
				if teacher.Name != expTeacher.Name {
					t.Errorf("unexpected teacher name at index %d: got %s, want %s", i, teacher.Name, expTeacher.Name)
				}
			}

			if len(lecture.LecturePlans) != len(tc.expected.LecturePlans) {
				t.Fatalf("unexpected number of lecture plans: got %d, want %d", len(lecture.LecturePlans), len(tc.expected.LecturePlans))
			}
			for i, plan := range lecture.LecturePlans {
				expPlan := tc.expected.LecturePlans[i]
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

			if len(lecture.Keywords) != len(tc.expected.Keywords) {
				t.Fatalf("unexpected number of keywords: got %d, want %d", len(lecture.Keywords), len(tc.expected.Keywords))
			}
			for i, keyword := range lecture.Keywords {
				expKeyword := tc.expected.Keywords[i]
				if keyword != expKeyword {
					t.Errorf("unexpected keyword at index %d: got %s, want %s", i, keyword, expKeyword)
				}
			}

			if len(lecture.RelatedCourseCodes) != len(tc.expected.RelatedCourseCodes) {
				t.Fatalf("unexpected number of related course codes: got %d, want %d", len(lecture.RelatedCourseCodes), len(tc.expected.RelatedCourseCodes))
			}
			for i, code := range lecture.RelatedCourseCodes {
				expCode := tc.expected.RelatedCourseCodes[i]
				if code != expCode {
					t.Errorf("unexpected related course code at index %d: got %s, want %s", i, code, expCode)
				}
			}

			if len(lecture.RelatedCourses) != len(tc.expected.RelatedCourses) {
				t.Fatalf("unexpected number of related courses: got %d, want %d", len(lecture.RelatedCourses), len(tc.expected.RelatedCourses))
			}
			for i, relatedCourse := range lecture.RelatedCourses {
				expRelatedCourse := tc.expected.RelatedCourses[i]
				if relatedCourse != expRelatedCourse {
					t.Errorf("unexpected related course at index %d: got %d, want %d", i, relatedCourse, expRelatedCourse)
				}
			}
		})
	}
}

func TestAddEnglishTitle(t *testing.T) {
	file, err := os.Open("fixture/course_detail_en.html")
	if err != nil {
		t.Fatalf("open course detail english fixture: %v", err)
	}
	defer file.Close()

	lecture := &domain.Lecture{}
	if err := AddEnglishTitle(file, lecture); err != nil {
		t.Fatalf("AddEnglishTitle returned error: %v", err)
	}
	if lecture.EnglishTitle != "Constitutional Law A" {
		t.Fatalf("unexpected english title: %s", lecture.EnglishTitle)
	}
}

func TestParseTimetablesRangeExpansion(t *testing.T) {
	testCases := []struct {
		name string
		raw  string
	}{
		{name: "hyphen", raw: "月5-8"},
		{name: "fullwidth wave", raw: "月5～8"},
		{name: "tilde", raw: "月5〜8"},
	}

	expectedPeriods := []domain.Period{domain.Period5, domain.Period6, domain.Period7, domain.Period8}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := parseTimetables(tc.raw, "3Q")
			if len(got) != len(expectedPeriods) {
				t.Fatalf("unexpected timetable count: got %d, want %d", len(got), len(expectedPeriods))
			}
			for i, tt := range got {
				if tt.Semester != domain.SemesterFall {
					t.Errorf("unexpected semester at index %d: got %s, want %s", i, tt.Semester, domain.SemesterFall)
				}
				if tt.DayOfWeek != domain.DayOfWeekMonday {
					t.Errorf("unexpected day at index %d: got %s, want %s", i, tt.DayOfWeek, domain.DayOfWeekMonday)
				}
				if tt.Period != expectedPeriods[i] {
					t.Errorf("unexpected period at index %d: got %d, want %d", i, tt.Period, expectedPeriods[i])
				}
			}
		})
	}
}
