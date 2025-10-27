package scraper

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const TopPageURL = "https://syllabus.s.isct.ac.jp"

// ListCoursesPagesURL fetches the top page and lists course URLs for the given year.
func ListCoursesPagesURL(year int) []string {
	resp, err := http.Get(TopPageURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	prefix := fmt.Sprintf("%s/courses/%d/", TopPageURL, year)
	pattern := regexp.MustCompile(regexp.QuoteMeta(prefix) + `[^"']*`)
	matches := pattern.FindAllString(string(body), -1)
	if len(matches) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{})
	urls := make([]string, 0, len(matches))

	for _, url := range matches {
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		urls = append(urls, url)
	}

	return urls
}
