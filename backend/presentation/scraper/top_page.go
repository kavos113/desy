package scraper

import (
	"fmt"
	"io"
	"regexp"
)

const TopPageURL = "https://syllabus.s.isct.ac.jp"

// ListCoursesPagesURL extracts course page URLs for the given year from the provided top page HTML.
func ListCoursesPagesURL(r io.Reader, year int) ([]string, error) {
	if r == nil {
		return nil, fmt.Errorf("nil reader provided")
	}

	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read top page html: %w", err)
	}

	prefix := fmt.Sprintf("%s/courses/%d/", TopPageURL, year)
	pattern := regexp.MustCompile(regexp.QuoteMeta(prefix) + `[^"']+`)
	matches := pattern.FindAllString(string(body), -1)
	if len(matches) == 0 {
		return []string{}, nil
	}

	seen := make(map[string]struct{}, len(matches))
	urls := make([]string, 0, len(matches))

	for _, url := range matches {
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		urls = append(urls, url)
	}

	return urls, nil
}
