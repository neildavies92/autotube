package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neildavies92/autotube/apps/api/internal/providers"
)

// fakeYouTube is a test double for providers.YouTubeProvider.
type fakeYouTube struct {
	searchResult providers.YouTubeSearchResult
	searchErr    error
	metricsResult providers.VideoMetricsResult
	metricsErr   error
}

func (f *fakeYouTube) SearchVideos(_ context.Context, _ providers.YouTubeSearchRequest) (providers.YouTubeSearchResult, error) {
	return f.searchResult, f.searchErr
}

func (f *fakeYouTube) GetVideoMetrics(_ context.Context, _ providers.VideoMetricsRequest) (providers.VideoMetricsResult, error) {
	return f.metricsResult, f.metricsErr
}

func TestRootRoute(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), nil).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json, got %q", got)
	}

	var body rootResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Service != "autotube-api" {
		t.Fatalf("expected service autotube-api, got %q", body.Service)
	}
}

func TestHealthRoute(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), nil).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body healthResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Service != "autotube-api" {
		t.Fatalf("expected service autotube-api, got %q", body.Service)
	}

	if body.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Status)
	}

	if body.CheckedAt == "" {
		t.Fatal("expected checkedAt to be set")
	}
}

func TestUnknownRoute(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), nil).ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestYouTubeResearchRoute(t *testing.T) {
	t.Parallel()

	yt := &fakeYouTube{
		searchResult: providers.YouTubeSearchResult{
			Videos: []providers.YouTubeVideo{
				{ID: "vid1", Title: "Finance Tips", ChannelName: "MoneyPro", PublishedAt: "2026-01-01T00:00:00Z"},
			},
		},
		metricsResult: providers.VideoMetricsResult{
			Metrics: []providers.VideoMetrics{
				{VideoID: "vid1", ViewCount: 50000, LikeCount: 1200, CommentCount: 300},
			},
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/research/youtube?q=personal+finance&region=GB", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), yt).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body researchResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Query != "personal finance" {
		t.Errorf("expected query %q, got %q", "personal finance", body.Query)
	}
	if body.Region != "GB" {
		t.Errorf("expected region GB, got %q", body.Region)
	}
	if body.Cached {
		t.Error("expected cached=false on first call")
	}
	if body.QuotaCost != 101 {
		t.Errorf("expected quota_cost 101, got %d", body.QuotaCost)
	}
	if len(body.Videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(body.Videos))
	}

	v := body.Videos[0]
	if v.VideoID != "vid1" {
		t.Errorf("expected video ID vid1, got %q", v.VideoID)
	}
	if v.ViewCount != 50000 {
		t.Errorf("expected 50000 views, got %d", v.ViewCount)
	}
	if v.VideoURL != "https://www.youtube.com/watch?v=vid1" {
		t.Errorf("unexpected video URL %q", v.VideoURL)
	}
}

func TestYouTubeResearchRouteMissingQuery(t *testing.T) {
	t.Parallel()

	yt := &fakeYouTube{}
	request := httptest.NewRequest(http.MethodGet, "/research/youtube", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), yt).ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestYouTubeResearchRouteNotRegisteredWithoutProvider(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/research/youtube?q=test", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), nil).ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d when no provider, got %d", http.StatusNotFound, response.Code)
	}
}

func TestYouTubeResearchRouteSearchError(t *testing.T) {
	t.Parallel()

	yt := &fakeYouTube{searchErr: errors.New("quota exceeded")}
	request := httptest.NewRequest(http.MethodGet, "/research/youtube?q=finance", nil)
	response := httptest.NewRecorder()

	New(slog.New(slog.DiscardHandler), yt).ServeHTTP(response, request)

	if response.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, response.Code)
	}
}

func TestYouTubeResearchRouteCache(t *testing.T) {
	t.Parallel()

	callCount := 0
	yt := &fakeYouTube{}
	yt.searchResult = providers.YouTubeSearchResult{
		Videos: []providers.YouTubeVideo{{ID: "vid1", Title: "Test"}},
	}

	// Wrap to count calls.
	counting := &countingYouTube{inner: yt, calls: &callCount}

	handler := New(slog.New(slog.DiscardHandler), counting)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/research/youtube?q=finance&region=US", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
	}

	if callCount != 1 {
		t.Errorf("expected SearchVideos to be called once (cache), got %d calls", callCount)
	}
}

type countingYouTube struct {
	inner *fakeYouTube
	calls *int
}

func (c *countingYouTube) SearchVideos(ctx context.Context, req providers.YouTubeSearchRequest) (providers.YouTubeSearchResult, error) {
	*c.calls++
	return c.inner.SearchVideos(ctx, req)
}

func (c *countingYouTube) GetVideoMetrics(ctx context.Context, req providers.VideoMetricsRequest) (providers.VideoMetricsResult, error) {
	return c.inner.GetVideoMetrics(ctx, req)
}
