package youtube

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/neildavies92/autotube/apps/api/internal/providers"
)

func newTestClient(srv *httptest.Server) *Client {
	return &Client{
		apiKey:     "test-key",
		baseURL:    srv.URL,
		httpClient: srv.Client(),
	}
}

func TestSearchVideos(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("key") == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(searchListResponse{
			Items: []searchItem{
				{
					ID:      struct{ VideoID string `json:"videoId"` }{VideoID: "vid1"},
					Snippet: struct {
						Title        string `json:"title"`
						Description  string `json:"description"`
						ChannelID    string `json:"channelId"`
						ChannelTitle string `json:"channelTitle"`
						PublishedAt  string `json:"publishedAt"`
					}{
						Title:        "Personal Finance Tips",
						Description:  "Learn how to save",
						ChannelID:    "ch1",
						ChannelTitle: "Finance Pro",
						PublishedAt:  "2026-01-01T00:00:00Z",
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	result, err := client.SearchVideos(context.Background(), providers.YouTubeSearchRequest{
		Query:      "personal finance",
		Locale:     providers.Locale{Region: "GB"},
		MaxResults: 5,
	})
	if err != nil {
		t.Fatalf("SearchVideos: %v", err)
	}

	if len(result.Videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(result.Videos))
	}

	video := result.Videos[0]
	if video.ID != "vid1" {
		t.Errorf("expected video ID vid1, got %q", video.ID)
	}
	if video.Title != "Personal Finance Tips" {
		t.Errorf("expected title %q, got %q", "Personal Finance Tips", video.Title)
	}
	if video.ChannelName != "Finance Pro" {
		t.Errorf("expected channel Finance Pro, got %q", video.ChannelName)
	}
}

func TestGetVideoMetrics(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/videos" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(videoListResponse{
			Items: []videoItem{
				{
					ID: "vid1",
					Statistics: struct {
						ViewCount    string `json:"viewCount"`
						LikeCount    string `json:"likeCount"`
						CommentCount string `json:"commentCount"`
					}{
						ViewCount:    "125000",
						LikeCount:    "3200",
						CommentCount: "450",
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	result, err := client.GetVideoMetrics(context.Background(), providers.VideoMetricsRequest{
		VideoIDs: []string{"vid1"},
	})
	if err != nil {
		t.Fatalf("GetVideoMetrics: %v", err)
	}

	if len(result.Metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(result.Metrics))
	}

	m := result.Metrics[0]
	if m.VideoID != "vid1" {
		t.Errorf("expected video ID vid1, got %q", m.VideoID)
	}
	if m.ViewCount != 125000 {
		t.Errorf("expected 125000 views, got %d", m.ViewCount)
	}
	if m.LikeCount != 3200 {
		t.Errorf("expected 3200 likes, got %d", m.LikeCount)
	}
}

func TestGetVideoMetricsEmpty(t *testing.T) {
	t.Parallel()

	// No HTTP calls should be made for an empty request.
	client := &Client{apiKey: "test", baseURL: "http://should-not-be-called", httpClient: http.DefaultClient}

	result, err := client.GetVideoMetrics(context.Background(), providers.VideoMetricsRequest{})
	if err != nil {
		t.Fatalf("expected no error for empty request, got %v", err)
	}
	if len(result.Metrics) != 0 {
		t.Fatalf("expected 0 metrics, got %d", len(result.Metrics))
	}
}

func TestSearchVideosAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(apiErrorResponse{})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	_, err := client.SearchVideos(context.Background(), providers.YouTubeSearchRequest{Query: "test"})
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}
