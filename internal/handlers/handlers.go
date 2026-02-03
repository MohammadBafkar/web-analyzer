package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/mohammadbafkar/web-analyzer/internal/analyzer"
	"github.com/mohammadbafkar/web-analyzer/internal/fetcher"
)

func HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func HandleReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func HandleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func HandleAnalyze(c *gin.Context) {
	url := strings.TrimSpace(c.PostForm("url"))

	if url == "" {
		renderError(c, http.StatusBadRequest, "Please enter a URL to analyze", 0)
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	result, err := fetcher.FetchURL(url)
	if err != nil {
		var fetchErr *fetcher.Error
		if errors.As(err, &fetchErr) {
			renderError(c, http.StatusUnprocessableEntity, fetchErr.Message, fetchErr.StatusCode)
		} else {
			renderError(c, http.StatusUnprocessableEntity, err.Error(), 0)
		}
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result.HTML))
	if err != nil {
		renderError(c, http.StatusUnprocessableEntity, "Failed to parse HTML: "+err.Error(), 0)
		return
	}

	analysis, err := analyzer.AnalyzePage(doc, result.HTML, result.FinalURL)
	if err != nil {
		renderError(c, http.StatusInternalServerError, "Analysis failed: "+err.Error(), 0)
		return
	}

	c.HTML(http.StatusOK, "results.html", gin.H{
		"Result": analysis,
	})
}

func renderError(c *gin.Context, httpStatus int, message string, remoteStatusCode int) {
	c.HTML(httpStatus, "error.html", gin.H{
		"Message":          message,
		"RemoteStatusCode": remoteStatusCode,
	})
}
