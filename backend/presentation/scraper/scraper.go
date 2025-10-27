package scraper

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kavos113/desy/backend/domain"
	"golang.org/x/net/html"
)

// CourseListItem represents a single row extracted from the course list page.
type CourseListItem struct {
	Code      string
	Title     string
	DetailURL string
}

// ParseCourseList extracts course metadata and detail links from a course list HTML document.
func ParseCourseList(r io.Reader, base string) ([]CourseListItem, error) {
	if r == nil {
		return nil, fmt.Errorf("nil reader provided")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("parse course list html: %w", err)
	}

	var resolvedBase *url.URL
	if base != "" {
		resolvedBase, err = url.Parse(base)
		if err != nil {
			return nil, fmt.Errorf("parse base url: %w", err)
		}
	}

	items := make([]CourseListItem, 0)
	doc.Find("table.c-table tbody tr").Each(func(_ int, tr *goquery.Selection) {
		cols := tr.Find("td")
		if cols.Length() < 2 {
			return
		}

		code := strings.TrimSpace(cols.Eq(0).Text())
		link := cols.Eq(1).Find("a")
		if code == "" || link.Length() == 0 {
			return
		}

		href, exists := link.Attr("href")
		if !exists || strings.TrimSpace(href) == "" {
			return
		}

		detail := strings.TrimSpace(href)
		if resolvedBase != nil {
			if u, err := resolvedBase.Parse(detail); err == nil {
				detail = u.String()
			}
		}

		title := strings.TrimSpace(normalizeWhitespace(selectionToText(link)))
		items = append(items, CourseListItem{Code: code, Title: title, DetailURL: detail})
	})

	return items, nil
}

// ParseCourseDetail scrapes a lecture aggregate from a course detail HTML document.
func ParseCourseDetail(r io.Reader, detailURL string) (*domain.Lecture, error) {
	if r == nil {
		return nil, fmt.Errorf("nil reader provided")
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("parse course detail html: %w", err)
	}

	lecture := &domain.Lecture{Url: strings.TrimSpace(detailURL)}

	lecture.Title = strings.TrimSpace(doc.Find("h1.c-h1").First().Text())
	lecture.Department = extractDefinition(doc, "開講元")
	lecture.Teachers = parseTeachers(extractDefinition(doc, "担当教員"))
	lecture.LectureType = parseLectureType(extractDefinition(doc, "授業形態"))
	lecture.Code = extractDefinition(doc, "科目コード")
	lecture.Credit = parseCredit(extractDefinition(doc, "単位数"))
	lecture.Year = parseFirstInt(extractDefinition(doc, "開講時期"))
	quarter := extractDefinition(doc, "開講クォーター")
	lecture.Language = extractDefinition(doc, "使用言語")

	lecture.Abstract = extractSectionText(doc, "授業の目的（ねらい）、概要")
	lecture.Goal = extractSectionText(doc, "到達目標")
	lecture.Experience = extractSectionText(doc, "実務経験のある教員等による授業科目等")
	lecture.Flow = extractSectionText(doc, "授業の進め方")
	lecture.OutOfClassWork = extractSectionText(doc, "準備学修(事前学修・復習)等についての指示")
	lecture.Textbook = extractSectionText(doc, "教科書")
	lecture.ReferenceBook = extractSectionText(doc, "参考書、講義資料等")
	lecture.Assessment = extractSectionText(doc, "成績評価の方法及び基準")
	lecture.Prerequisite = extractSectionText(doc, "履修の条件・注意事項")
	lecture.Note = extractSectionText(doc, "その他")

	lecture.Keywords = parseKeywords(extractSectionText(doc, "キーワード"))
	lecture.LecturePlans = parseLecturePlans(doc)
	lecture.Timetables = parseTimetables(extractDefinitionRaw(doc, "曜日・時限"), quarter)

	return lecture, nil
}

func extractDefinition(doc *goquery.Document, term string) string {
	sel := extractDefinitionSelection(doc, term)
	if sel == nil || sel.Length() == 0 {
		return ""
	}
	return normalizeWhitespace(selectionToText(sel))
}

func extractDefinitionRaw(doc *goquery.Document, term string) string {
	sel := extractDefinitionSelection(doc, term)
	if sel == nil || sel.Length() == 0 {
		return ""
	}
	html, err := sel.Html()
	if err != nil {
		return normalizeWhitespace(sel.Text())
	}
	return normalizeWhitespace(strings.ReplaceAll(strings.ReplaceAll(html, "<br>", "\n"), "<br />", "\n"))
}

func extractDefinitionSelection(doc *goquery.Document, term string) *goquery.Selection {
	var result *goquery.Selection
	doc.Find("div.c-dl-2col__item").EachWithBreak(func(_ int, item *goquery.Selection) bool {
		dt := normalizeWhitespace(item.Find("dt").First().Text())
		if dt == term || strings.Contains(dt, term) {
			result = item.Find("dd").First()
			return false
		}
		return true
	})
	if result == nil {
		return nil
	}
	return result
}

func extractSectionText(doc *goquery.Document, heading string) string {
	var text string
	doc.Find("h3.c-h3").EachWithBreak(func(_ int, h3 *goquery.Selection) bool {
		if strings.TrimSpace(h3.Text()) == heading {
			var sb strings.Builder
			for s := h3.Next(); s.Length() > 0; s = s.Next() {
				if goquery.NodeName(s) == "h3" {
					break
				}
				sb.WriteString(selectionToText(s))
				sb.WriteString("\n")
			}
			text = normalizeWhitespace(sb.String())
			return false
		}
		return true
	})
	return text
}

func selectionToText(sel *goquery.Selection) string {
	if sel == nil || sel.Length() == 0 {
		return ""
	}
	var sb strings.Builder
	sel.Each(func(_ int, s *goquery.Selection) {
		if len(s.Nodes) == 0 {
			return
		}
		nodeToText(s.Nodes[0], &sb)
	})
	return normalizeWhitespace(sb.String())
}

func nodeToText(n *html.Node, sb *strings.Builder) {
	if n == nil {
		return
	}
	switch n.Type {
	case html.TextNode:
		sb.WriteString(n.Data)
	case html.ElementNode:
		name := strings.ToLower(n.Data)
		switch name {
		case "br":
			sb.WriteString("\n")
		case "li", "p", "div", "tr":
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				nodeToText(c, sb)
			}
			sb.WriteString("\n")
			return
		default:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				nodeToText(c, sb)
			}
		}
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nodeToText(c, sb)
		}
	}
}

var whitespaceReducer = regexp.MustCompile(`(?m)[ \t\f\v]+`)

func normalizeWhitespace(input string) string {
	if input == "" {
		return ""
	}
	cleaned := strings.ReplaceAll(strings.ReplaceAll(input, "\r\n", "\n"), "\r", "\n")
	cleaned = whitespaceReducer.ReplaceAllString(cleaned, " ")
	lines := strings.Split(cleaned, "\n")
	result := make([]string, 0, len(lines))
	lastEmpty := true
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !lastEmpty {
				result = append(result, "")
			}
			lastEmpty = true
			continue
		}
		lastEmpty = false
		result = append(result, trimmed)
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

func parseTeachers(raw string) []domain.Teacher {
	if raw == "" {
		return nil
	}
	parts := splitList(raw)
	teachers := make([]domain.Teacher, 0, len(parts))
	for _, name := range parts {
		if name == "" {
			continue
		}
		teachers = append(teachers, domain.Teacher{Name: name})
	}
	return teachers
}

func parseLectureType(raw string) domain.LectureType {
	lower := strings.ToLower(raw)
	switch {
	case strings.Contains(lower, "対面") || strings.Contains(lower, "教室"):
		return domain.LectureTypeOffline
	case strings.Contains(lower, "ライブ") || strings.Contains(lower, "リアルタイム"):
		return domain.LectureTypeLive
	case strings.Contains(lower, "ハイフレックス") || strings.Contains(lower, "ハイブリッド"):
		return domain.LectureTypeHyflex
	case strings.Contains(lower, "オンデマンド") || strings.Contains(lower, "録画"):
		return domain.LectureTypeOndemand
	case strings.TrimSpace(raw) == "":
		return ""
	default:
		return domain.LectureTypeOther
	}
}

var numberRegexp = regexp.MustCompile(`(\d+)`)

func parseFirstInt(raw string) int {
	if matches := numberRegexp.FindStringSubmatch(raw); len(matches) > 0 {
		value, err := strconv.Atoi(matches[1])
		if err == nil {
			return value
		}
	}
	return 0
}

func parseCredit(raw string) int {
	value := parseFirstInt(raw)
	switch {
	case value >= 100 && value%100 == 0:
		return value / 100
	case value >= 10 && value%10 == 0:
		return value / 10
	default:
		return value
	}
}

func parseKeywords(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := splitList(raw)
	keywords := make([]string, 0, len(parts))
	for _, keyword := range parts {
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func splitList(raw string) []string {
	separators := []string{"/", "／", "・", ",", "，", ";", "|"}
	normalized := strings.ReplaceAll(strings.ReplaceAll(raw, "\n", ","), "\r", ",")
	normalized = strings.ReplaceAll(normalized, "　", " ")
	for _, sep := range []string{"、", ";", "|"} {
		normalized = strings.ReplaceAll(normalized, sep, ",")
	}
	for _, sep := range separators {
		normalized = strings.ReplaceAll(normalized, sep, ",")
	}
	tokens := strings.Split(normalized, ",")
	results := make([]string, 0, len(tokens))
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}

func parseLecturePlans(doc *goquery.Document) []domain.LecturePlan {
	plans := make([]domain.LecturePlan, 0)
	doc.Find("table#lecture_plans tbody tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 3 {
			return
		}
		count := parseFirstInt(cells.Eq(0).Text())
		plan := selectionToText(cells.Eq(1))
		assignment := selectionToText(cells.Eq(2))
		if count == 0 && plan == "" && assignment == "" {
			return
		}
		plans = append(plans, domain.LecturePlan{Count: count, Plan: plan, Assignment: assignment})
	})
	return plans
}

var dayOfWeekMap = map[rune]domain.DayOfWeek{
	'月': domain.DayOfWeekMonday,
	'火': domain.DayOfWeekTuesday,
	'水': domain.DayOfWeekWednesday,
	'木': domain.DayOfWeekThursday,
	'金': domain.DayOfWeekFriday,
	'土': domain.DayOfWeekSaturday,
	'日': domain.DayOfWeekSunday,
}

func parseTimetables(raw, quarter string) []domain.TimeTable {
	raw = normalizeWhitespace(raw)
	if raw == "" || raw == "-" {
		return nil
	}
	entries := splitLines(raw)
	timetables := make([]domain.TimeTable, 0, len(entries))
	for _, entry := range entries {
		if entry == "" {
			continue
		}
		var roomName string
		if idx := strings.Index(entry, "("); idx != -1 {
			if end := strings.LastIndex(entry, ")"); end > idx {
				roomName = strings.TrimSpace(entry[idx+1 : end])
				entry = strings.TrimSpace(entry[:idx])
			}
		}
		if entry == "" {
			continue
		}
		dayRune := []rune(entry)[0]
		day, ok := dayOfWeekMap[dayRune]
		if !ok {
			continue
		}
		periodPart := strings.TrimSpace(entry[1:])
		numbers := numberRegexp.FindAllString(periodPart, -1)
		if len(numbers) == 0 {
			tt := domain.TimeTable{Semester: domain.Semester(quarter), DayOfWeek: day}
			if roomName != "" {
				tt.Room.Name = roomName
			}
			timetables = append(timetables, tt)
			continue
		}
		for _, num := range numbers {
			p, err := strconv.Atoi(num)
			if err != nil {
				continue
			}
			tt := domain.TimeTable{Semester: domain.Semester(quarter), DayOfWeek: day, Period: domain.Period(p)}
			if roomName != "" {
				tt.Room.Name = roomName
			}
			timetables = append(timetables, tt)
		}
	}
	return timetables
}

func splitLines(raw string) []string {
	if raw == "" {
		return nil
	}
	lines := strings.Split(raw, "\n")
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}
