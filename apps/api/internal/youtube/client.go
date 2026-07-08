package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/neildavies92/autotube/apps/api/internal/providers"
)

const (
	defaultBaseURL    = "https://www.googleapis.com/youtube/v3"
	defaultMaxResults = 25
)

// Client calls the YouTube Data API v3.
// Search costs 100 quota units; video metrics cost 1 unit per call.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: httpClient,
	}
}

func (c *Client) SearchVideos(ctx context.Context, req providers.YouTubeSearchRequest) (providers.YouTubeSearchResult, error) {
	maxResults := req.MaxResults
	if maxResults <= 0 {
		maxResults = defaultMaxResults
	}

	params := url.Values{
		"part":       {"snippet"},
		"q":          {req.Query},
		"type":       {"video"},
		"maxResults": {strconv.Itoa(maxResults)},
		"key":        {c.apiKey},
	}
	if req.Locale.Region != "" {
		params.Set("regionCode", req.Locale.Region)
	}
	if req.Locale.Language != "" {
		params.Set("relevanceLanguage", req.Locale.Language)
	}

	var searchResp searchListResponse
	if err := c.get(ctx, "/search", params, &searchResp); err != nil {
		return providers.YouTubeSearchResult{}, fmt.Errorf("youtube: search: %w", err)
	}

	videos := make([]providers.YouTubeVideo, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		videos = append(videos, providers.YouTubeVideo{
			ID:          item.ID.VideoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			ChannelID:   item.Snippet.ChannelID,
			ChannelName: item.Snippet.ChannelTitle,
			PublishedAt: item.Snippet.PublishedAt,
		})
	}

	return providers.YouTubeSearchResult{Videos: videos}, nil
}

func (c *Client) GetVideoMetrics(ctx context.Context, req providers.VideoMetricsRequest) (providers.VideoMetricsResult, error) {
	if len(req.VideoIDs) == 0 {
		return providers.VideoMetricsResult{}, nil
	}

	params := url.Values{
		"part": {"statistics"},
		"id":   {strings.Join(req.VideoIDs, ",")},
		"key":  {c.apiKey},
	}

	var videosResp videoListResponse
	if err := c.get(ctx, "/videos", params, &videosResp); err != nil {
		return providers.VideoMetricsResult{}, fmt.Errorf("youtube: video metrics: %w", err)
	}

	metrics := make([]providers.VideoMetrics, 0, len(videosResp.Items))
	for _, item := range videosResp.Items {
		metrics = append(metrics, providers.VideoMetrics{
			VideoID:      item.ID,
			ViewCount:    parseCount(item.Statistics.ViewCount),
			LikeCount:    parseCount(item.Statistics.LikeCount),
			CommentCount: parseCount(item.Statistics.CommentCount),
		})
	}

	return providers.VideoMetricsResult{Metrics: metrics}, nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dest any) error {
	endpoint := c.baseURL + path + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr apiErrorResponse
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr == nil && apiErr.Error.Message != "" {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, apiErr.Error.Message)
		}
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// YouTube Data API v3 response shapes.

type apiErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type searchListResponse struct {
	Items []searchItem `json:"items"`
}

type searchItem struct {
	ID struct {
		VideoID string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		ChannelID    string `json:"channelId"`
		ChannelTitle string `json:"channelTitle"`
		PublishedAt  string `json:"publishedAt"`
	} `json:"snippet"`
}

type videoListResponse struct {
	Items []videoItem `json:"items"`
}

type videoItem struct {
	ID         string `json:"id"`
	Statistics struct {
		ViewCount    string `json:"viewCount"`
		LikeCount    string `json:"likeCount"`
		CommentCount string `json:"commentCount"`
	} `json:"statistics"`
}

func parseCount(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
