package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
	}
	return e.Message
}

type Result struct {
	HTML       string
	StatusCode int
	FinalURL   string
}

func FetchURL(url string) (*Result, error) {
	if url == "" {
		return nil, &Error{Message: "URL cannot be empty"}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return &Error{Message: "too many redirects (max 10)"}
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("invalid URL: %v", err)}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebAnalyzer/1.0)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to fetch URL: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Message:    getStatusMessage(resp.StatusCode),
		}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to read response: %v", err)}
	}

	return &Result{
		HTML:       string(body),
		StatusCode: resp.StatusCode,
		FinalURL:   resp.Request.URL.String(),
	}, nil
}

func ExtractHost(urlStr string) string {
	re := regexp.MustCompile(`^(?:https?://)?([^/]+)`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}

func ResolveURL(baseURL, href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	if strings.HasPrefix(href, "//") {
		if strings.HasPrefix(baseURL, "https://") {
			return "https:" + href
		}
		return "http:" + href
	}

	baseRe := regexp.MustCompile(`^(https?://[^/]+)`)
	baseMatches := baseRe.FindStringSubmatch(baseURL)
	if len(baseMatches) < 2 {
		return ""
	}
	base := baseMatches[1]

	if strings.HasPrefix(href, "/") {
		return base + href
	}

	pathRe := regexp.MustCompile(`^(https?://[^/]+(?:/[^?#]*)?)`)
	pathMatches := pathRe.FindStringSubmatch(baseURL)
	if len(pathMatches) < 2 {
		return base + "/" + href
	}

	path := pathMatches[1]
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash > len(base) {
		path = path[:lastSlash+1]
	} else {
		path = base + "/"
	}

	return path + href
}

func CheckLinkAccessibility(links []string) int {
	if len(links) == 0 {
		return 0
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	inaccessible := 0
	sem := make(chan struct{}, 10)
	for _, link := range links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			req, err := http.NewRequest("HEAD", url, nil)
			if err != nil {
				mu.Lock()
				inaccessible++
				mu.Unlock()
				return
			}

			req.Header.Set("User-Agent", "WebAnalyzer/1.0")

			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				inaccessible++
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				mu.Lock()
				inaccessible++
				mu.Unlock()
			}
		}(link)
	}

	wg.Wait()
	return inaccessible
}

func getStatusMessage(code int) string {
	messages := map[int]string{
		400: "Bad Request - The server could not understand the request",
		401: "Unauthorized - Authentication is required",
		403: "Forbidden - Access to this resource is denied",
		404: "Not Found - The requested page does not exist",
		405: "Method Not Allowed - The request method is not supported",
		408: "Request Timeout - The server timed out waiting for the request",
		429: "Too Many Requests - Rate limit exceeded",
		500: "Internal Server Error - The server encountered an error",
		502: "Bad Gateway - The server received an invalid response",
		503: "Service Unavailable - The server is temporarily unavailable",
		504: "Gateway Timeout - The server did not respond in time",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}
	return http.StatusText(code)
}
