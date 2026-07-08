package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/neildavies92/autotube/apps/api/internal/providers"
)

const researchCacheTTL = 4 * time.Hour

type server struct {
	logger        *slog.Logger
	youtube       providers.YouTubeProvider
	researchCache *researchCache
}

// New builds the HTTP handler. Pass a non-nil YouTubeProvider to enable the
// /research/youtube endpoint; pass nil to omit it (e.g. no API key configured).
func New(logger *slog.Logger, yt providers.YouTubeProvider) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	s := &server{
		logger:        logger,
		youtube:       yt,
		researchCache: newResearchCache(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleRoot)
	mux.HandleFunc("GET /health", handleHealth)

	if yt != nil {
		mux.HandleFunc("GET /research/youtube", s.handleYouTubeResearch)
	}

	return recoverer(logger, requestLogger(logger, mux))
}

// --- static handlers ---

type rootResponse struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

type healthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CheckedAt string `json:"checkedAt"`
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, rootResponse{
		Service: "autotube-api",
		Message: "AutoTube API scaffold is running.",
	})
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Service:   "autotube-api",
		Status:    "ok",
		Message:   "AutoTube API is ready for local dashboard requests.",
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// --- YouTube research handler ---

type researchResponse struct {
	Query     string          `json:"query"`
	Region    string          `json:"region"`
	Cached    bool            `json:"cached"`
	FetchedAt string          `json:"fetchedAt"`
	QuotaCost int             `json:"quotaCost"`
	Videos    []researchVideo `json:"videos"`
}

type researchVideo struct {
	VideoID      string `json:"videoId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelID    string `json:"channelId"`
	ChannelName  string `json:"channelName"`
	PublishedAt  string `json:"publishedAt"`
	ViewCount    int64  `json:"viewCount"`
	LikeCount    int64  `json:"likeCount"`
	CommentCount int64  `json:"commentCount"`
	VideoURL     string `json:"videoUrl"`
}

func (s *server) handleYouTubeResearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, `{"error":"q is required"}`, http.StatusBadRequest)
		return
	}

	region := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("region")))
	if region == "" {
		region = "US"
	}

	maxResults := 25
	if maxParam := r.URL.Query().Get("max"); maxParam != "" {
		if n, err := strconv.Atoi(maxParam); err == nil && n > 0 && n <= 50 {
			maxResults = n
		}
	}

	cacheKey := fmt.Sprintf("%s|%s|%d", strings.ToLower(query), region, maxResults)

	if entry, ok := s.researchCache.get(cacheKey); ok {
		s.logger.Info("youtube research cache hit", "query", query, "region", region)
		writeJSON(w, http.StatusOK, researchResponse{
			Query:     query,
			Region:    region,
			Cached:    true,
			FetchedAt: entry.fetchedAt.UTC().Format(time.RFC3339),
			QuotaCost: 0,
			Videos:    entry.videos,
		})
		return
	}

	locale := providers.Locale{Region: region}

	searchResult, err := s.youtube.SearchVideos(r.Context(), providers.YouTubeSearchRequest{
		Query:      query,
		Locale:     locale,
		MaxResults: maxResults,
	})
	if err != nil {
		s.logger.Error("youtube search failed", "query", query, "error", err)
		http.Error(w, `{"error":"YouTube search failed: `+err.Error()+`"}`, http.StatusBadGateway)
		return
	}

	videoIDs := make([]string, 0, len(searchResult.Videos))
	for _, v := range searchResult.Videos {
		videoIDs = append(videoIDs, v.ID)
	}

	metricsResult, err := s.youtube.GetVideoMetrics(r.Context(), providers.VideoMetricsRequest{VideoIDs: videoIDs})
	if err != nil {
		s.logger.Error("youtube metrics failed", "query", query, "error", err)
		http.Error(w, `{"error":"YouTube metrics failed: `+err.Error()+`"}`, http.StatusBadGateway)
		return
	}

	metricsByID := make(map[string]providers.VideoMetrics, len(metricsResult.Metrics))
	for _, m := range metricsResult.Metrics {
		metricsByID[m.VideoID] = m
	}

	videos := make([]researchVideo, 0, len(searchResult.Videos))
	for _, v := range searchResult.Videos {
		m := metricsByID[v.ID]
		videos = append(videos, researchVideo{
			VideoID:      v.ID,
			Title:        v.Title,
			Description:  v.Description,
			ChannelID:    v.ChannelID,
			ChannelName:  v.ChannelName,
			PublishedAt:  v.PublishedAt,
			ViewCount:    m.ViewCount,
			LikeCount:    m.LikeCount,
			CommentCount: m.CommentCount,
			VideoURL:     "https://www.youtube.com/watch?v=" + v.ID,
		})
	}

	fetchedAt := time.Now()
	s.researchCache.set(cacheKey, researchCacheEntry{videos: videos, fetchedAt: fetchedAt}, researchCacheTTL)

	s.logger.Info("youtube research completed", "query", query, "region", region, "results", len(videos), "quota_cost", 101)

	writeJSON(w, http.StatusOK, researchResponse{
		Query:     query,
		Region:    region,
		Cached:    false,
		FetchedAt: fetchedAt.UTC().Format(time.RFC3339),
		QuotaCost: 101,
		Videos:    videos,
	})
}

// --- in-memory research cache ---

type researchCacheEntry struct {
	videos    []researchVideo
	fetchedAt time.Time
	expiresAt time.Time
}

type researchCache struct {
	mu      sync.RWMutex
	entries map[string]researchCacheEntry
}

func newResearchCache() *researchCache {
	return &researchCache{entries: make(map[string]researchCacheEntry)}
}

func (c *researchCache) get(key string) (researchCacheEntry, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		return researchCacheEntry{}, false
	}
	return entry, true
}

func (c *researchCache) set(key string, entry researchCacheEntry, ttl time.Duration) {
	entry.expiresAt = time.Now().Add(ttl)
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}

// --- middleware ---

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}

func recoverer(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered",
					"method", r.Method,
					"path", r.URL.Path,
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
