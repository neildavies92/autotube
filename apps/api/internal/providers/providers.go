package providers

import "context"

type Locale struct {
	Language string
	Region   string
}

type YouTubeProvider interface {
	SearchVideos(context.Context, YouTubeSearchRequest) (YouTubeSearchResult, error)
	GetVideoMetrics(context.Context, VideoMetricsRequest) (VideoMetricsResult, error)
}

type YouTubeSearchRequest struct {
	Query      string
	Locale     Locale
	MaxResults int
}

type YouTubeSearchResult struct {
	Videos []YouTubeVideo
}

type YouTubeVideo struct {
	ID          string
	Title       string
	Description string
	ChannelID   string
	ChannelName string
	PublishedAt string
}

type VideoMetricsRequest struct {
	VideoIDs []string
}

type VideoMetricsResult struct {
	Metrics []VideoMetrics
}

type VideoMetrics struct {
	VideoID      string
	ViewCount    int64
	LikeCount    int64
	CommentCount int64
}

type TrendsProvider interface {
	InterestOverTime(context.Context, TrendInterestRequest) (TrendInterestResult, error)
	RelatedQueries(context.Context, RelatedQueriesRequest) (RelatedQueriesResult, error)
}

type TrendInterestRequest struct {
	Keywords []string
	Locale   Locale
	Window   string
}

type TrendInterestResult struct {
	Series []TrendPoint
}

type TrendPoint struct {
	Keyword string
	Period  string
	Score   int
}

type RelatedQueriesRequest struct {
	Keyword string
	Locale  Locale
}

type RelatedQueriesResult struct {
	Queries []RelatedQuery
}

type RelatedQuery struct {
	Query string
	Score int
}

type LLMProvider interface {
	GenerateText(context.Context, TextGenerationRequest) (TextGenerationResult, error)
}

type TextGenerationRequest struct {
	SystemPrompt string
	UserPrompt   string
	Metadata     map[string]string
}

type TextGenerationResult struct {
	Text       string
	ProviderID string
}

type ImageProvider interface {
	GenerateImage(context.Context, ImageGenerationRequest) (ImageAsset, error)
}

type ImageGenerationRequest struct {
	Prompt      string
	AspectRatio string
	Metadata    map[string]string
}

type ImageAsset struct {
	URI        string
	ProviderID string
	MimeType   string
}

type VoiceProvider interface {
	SynthesizeSpeech(context.Context, SpeechSynthesisRequest) (AudioAsset, error)
}

type SpeechSynthesisRequest struct {
	Text      string
	VoiceName string
	Locale    Locale
	Metadata  map[string]string
}

type AudioAsset struct {
	URI        string
	ProviderID string
	MimeType   string
}

type RenderProvider interface {
	RenderVideo(context.Context, RenderVideoRequest) (VideoAsset, error)
}

type RenderVideoRequest struct {
	ScriptURI string
	AssetURIs []string
	Metadata  map[string]string
}

type VideoAsset struct {
	URI        string
	ProviderID string
	MimeType   string
}

type UploadProvider interface {
	UploadPrivateVideo(context.Context, UploadVideoRequest) (UploadVideoResult, error)
}

type UploadVideoRequest struct {
	VideoURI    string
	Title       string
	Description string
	Tags        []string
	Metadata    map[string]string
}

type UploadVideoResult struct {
	ProviderID string
	VideoID    string
	PrivateURL string
}
