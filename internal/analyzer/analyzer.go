package analyzer

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mohammadbafkar/web-analyzer/internal/fetcher"
	"golang.org/x/net/html"
)

type Result struct {
	URL               string
	HTMLVersion       string
	Title             string
	Headings          map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	HasLoginForm      bool
	AnalyzedAt        time.Time
}

func AnalyzePage(doc *goquery.Document, rawHTML string, baseURL string) (*Result, error) {
	result := &Result{
		URL:        baseURL,
		Headings:   make(map[string]int),
		AnalyzedAt: time.Now(),
	}

	result.HTMLVersion = detectHTMLVersion(rawHTML)

	result.Title = doc.Find("title").First().Text()
	result.Title = strings.TrimSpace(result.Title)

	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		count := doc.Find(tag).Length()
		if count > 0 {
			result.Headings[tag] = count
		}
	}

	result.InternalLinks, result.ExternalLinks, result.InaccessibleLinks = analyzeLinks(doc, baseURL)

	result.HasLoginForm = detectLoginForm(doc)

	return result, nil
}

func detectHTMLVersion(rawHTML string) string {
	reader := strings.NewReader(rawHTML)
	tokenizer := html.NewTokenizer(reader)

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.DoctypeToken {
			doctype := strings.ToLower(string(tokenizer.Raw()))

			if strings.Contains(doctype, "html") && !strings.Contains(doctype, "public") && !strings.Contains(doctype, "dtd") {
				return "HTML5"
			}
			if strings.Contains(doctype, "xhtml 1.0 strict") {
				return "XHTML 1.0 Strict"
			}
			if strings.Contains(doctype, "xhtml 1.0 transitional") {
				return "XHTML 1.0 Transitional"
			}
			if strings.Contains(doctype, "xhtml 1.1") {
				return "XHTML 1.1"
			}
			if strings.Contains(doctype, "html 4.01 strict") {
				return "HTML 4.01 Strict"
			}
			if strings.Contains(doctype, "html 4.01 transitional") {
				return "HTML 4.01 Transitional"
			}
			if strings.Contains(doctype, "html 4.01 frameset") {
				return "HTML 4.01 Frameset"
			}
			if strings.Contains(doctype, "html 4.01") {
				return "HTML 4.01"
			}
			return "Unknown (DOCTYPE found)"
		}
		if tt == html.StartTagToken {
			break
		}
	}

	return "Unknown (No DOCTYPE)"
}

func analyzeLinks(doc *goquery.Document, baseURL string) (internal, external, inaccessible int) {
	baseHost := fetcher.ExtractHost(baseURL)

	var links []string
	linkSet := make(map[string]bool)

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		if strings.HasPrefix(href, "#") ||
			strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "mailto:") ||
			strings.HasPrefix(href, "tel:") {
			return
		}

		fullURL := fetcher.ResolveURL(baseURL, href)
		if fullURL == "" {
			return
		}

		linkHost := fetcher.ExtractHost(fullURL)
		if linkHost == baseHost {
			internal++
		} else {
			external++
		}

		if !linkSet[fullURL] {
			linkSet[fullURL] = true
			links = append(links, fullURL)
		}
	})

	inaccessible = fetcher.CheckLinkAccessibility(links)

	return
}

func detectLoginForm(doc *goquery.Document) bool {
	if doc.Find("input[type='password'], input[type=\"password\"]").Length() > 0 {
		return true
	}

	loginPatterns := []string{"login", "signin", "sign-in", "auth", "log-in"}

	found := false
	doc.Find("form").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		action, _ := s.Attr("action")
		id, _ := s.Attr("id")
		class, _ := s.Attr("class")
		name, _ := s.Attr("name")

		attrs := strings.ToLower(action + " " + id + " " + class + " " + name)

		for _, pattern := range loginPatterns {
			if strings.Contains(attrs, pattern) {
				found = true
				return false
			}
		}
		return true
	})

	return found
}
