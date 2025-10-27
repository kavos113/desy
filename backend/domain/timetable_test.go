package domain

import "testing"

func TestSemesterFromQuarter(t *testing.T) {
	type testCase struct {
		quarter  string
		expected []Semester
	}
	tests := []testCase{
		{"1Q", []Semester{SemesterSpring}},
		{"2Q", []Semester{SemesterSummer}},
		{"3Q", []Semester{SemesterFall}},
		{"4Q", []Semester{SemesterWinter}},
		{"1-2Q", []Semester{SemesterSpring, SemesterSummer}},
		{"2-3Q", []Semester{SemesterSummer, SemesterFall}},
		{"3-4Q", []Semester{SemesterFall, SemesterWinter}},
		{"1-4Q", []Semester{SemesterSpring, SemesterSummer, SemesterFall, SemesterWinter}},
	}

	for _, tt := range tests {
		t.Run(tt.quarter, func(t *testing.T) {
			got := FromQuarter(tt.quarter)
			if len(got) != len(tt.expected) {
				t.Errorf("unexpected length: got %d, want %d", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("unexpected semester at index %d: got %s, want %s", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
