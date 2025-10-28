package scraper

import (
	"io"

	"github.com/kavos113/desy/backend/domain"
)

// Parser exposes parsing operations over syllabus HTML pages.
type Parser interface {
	ParseCourseList(r io.Reader, base string) ([]CourseListItem, error)
	ParseCourseDetail(r io.Reader, detailURL string) (*domain.Lecture, error)
	ListCoursesPagesURL(r io.Reader, year int) ([]string, error)
	AddEnglishTitle(r io.Reader, lecture *domain.Lecture) error
}

// NewParser returns the default parser implementation.
func NewParser() Parser {
	return htmlParser{}
}

type htmlParser struct{}

func (htmlParser) ParseCourseList(r io.Reader, base string) ([]CourseListItem, error) {
	return ParseCourseList(r, base)
}

func (htmlParser) ParseCourseDetail(r io.Reader, detailURL string) (*domain.Lecture, error) {
	return ParseCourseDetail(r, detailURL)
}

func (htmlParser) ListCoursesPagesURL(r io.Reader, year int) ([]string, error) {
	return ListCoursesPagesURL(r, year)
}

func (htmlParser) AddEnglishTitle(r io.Reader, lecture *domain.Lecture) error {
	return AddEnglishTitle(r, lecture)
}
