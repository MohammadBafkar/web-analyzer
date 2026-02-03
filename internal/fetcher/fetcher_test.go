package fetcher

import "testing"

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		href     string
		expected string
	}{
		{
			name:     "Absolute URL",
			baseURL:  "https://example.com/page",
			href:     "https://other.com/path",
			expected: "https://other.com/path",
		},
		{
			name:     "Root relative",
			baseURL:  "https://example.com/page",
			href:     "/other",
			expected: "https://example.com/other",
		},
		{
			name:     "Protocol relative",
			baseURL:  "https://example.com/page",
			href:     "//cdn.example.com/file",
			expected: "https://cdn.example.com/file",
		},
		{
			name:     "Relative path",
			baseURL:  "https://example.com/dir/page",
			href:     "other.html",
			expected: "https://example.com/dir/other.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveURL(tt.baseURL, tt.href)
			if result != tt.expected {
				t.Errorf("ResolveURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractHost(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/path", "example.com"},
		{"http://www.test.org:8080/page", "www.test.org:8080"},
		{"https://sub.domain.com", "sub.domain.com"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := ExtractHost(tt.url)
			if result != tt.expected {
				t.Errorf("ExtractHost(%v) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}
