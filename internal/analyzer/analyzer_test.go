package analyzer

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5",
			html:     "<!DOCTYPE html><html><head></head><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01 Strict",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd"><html></html>`,
			expected: "HTML 4.01",
		},
		{
			name:     "XHTML 1.0 Strict",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html></html>`,
			expected: "XHTML 1.0 Strict",
		},
		{
			name:     "No DOCTYPE",
			html:     "<html><head></head><body></body></html>",
			expected: "Unknown (No DOCTYPE)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectHTMLVersion(tt.html)
			if result != tt.expected {
				t.Errorf("detectHTMLVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAnalyzePage(t *testing.T) {
	html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test Page</title></head>
		<body>
			<h1>Main Heading</h1>
			<h2>Subheading 1</h2>
			<h2>Subheading 2</h2>
			<h3>Sub-subheading</h3>
			<a href="/internal">Internal Link</a>
			<a href="https://example.com">External Link</a>
		</body>
		</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	result, err := AnalyzePage(doc, html, "https://test.com")
	if err != nil {
		t.Fatalf("AnalyzePage() error = %v", err)
	}

	if result.Title != "Test Page" {
		t.Errorf("Title = %v, want %v", result.Title, "Test Page")
	}

	if result.HTMLVersion != "HTML5" {
		t.Errorf("HTMLVersion = %v, want %v", result.HTMLVersion, "HTML5")
	}

	if result.Headings["h1"] != 1 {
		t.Errorf("Headings[h1] = %v, want %v", result.Headings["h1"], 1)
	}
	if result.Headings["h2"] != 2 {
		t.Errorf("Headings[h2] = %v, want %v", result.Headings["h2"], 2)
	}
}

func TestDetectLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "Has password input",
			html:     `<html><body><form><input type="password" name="password"></form></body></html>`,
			expected: true,
		},
		{
			name:     "Has login class",
			html:     `<html><body><form class="login-form"><input type="text"></form></body></html>`,
			expected: true,
		},
		{
			name:     "No login form",
			html:     `<html><body><form><input type="text" name="search"><button>Search</button></form></body></html>`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			result := detectLoginForm(doc)
			if result != tt.expected {
				t.Errorf("detectLoginForm() = %v, want %v", result, tt.expected)
			}
		})
	}
}
