package scraper

import (
	"os"
	"testing"
)

func TestListCoursesPagesURLFromFixture(t *testing.T) {
	file, err := os.Open("fixture/top_page.html")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	urls, err := ListCoursesPagesURL(file, 2025)
	if err != nil {
		t.Fatalf("ListCoursesPagesURL returned error: %v", err)
	}

	if len(urls) == 0 {
		t.Fatalf("expected URLs, got none")
	}

	expects := map[string]bool{
		"https://syllabus.s.isct.ac.jp/courses/2025/4/0-904-342300-0-0": false,
		"https://syllabus.s.isct.ac.jp/courses/2025/7/0-907-0-110100-0": false,
	}

	unique := make(map[string]struct{})
	for _, url := range urls {
		unique[url] = struct{}{}
		if _, ok := expects[url]; ok {
			expects[url] = true
		}
	}

	if len(unique) != len(urls) {
		t.Fatalf("expected unique URLs, got duplicates")
	}

	for target, seen := range expects {
		if !seen {
			t.Fatalf("expected URL %s not found", target)
		}
	}
}

func TestListCoursesPagesURLNilReader(t *testing.T) {
	if _, err := ListCoursesPagesURL(nil, 2025); err == nil {
		t.Fatalf("expected error for nil reader")
	}
}
