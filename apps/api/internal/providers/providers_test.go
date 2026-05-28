package providers

import (
	"context"
	"testing"
)

func TestProviderInterfacesAreComposable(t *testing.T) {
	t.Parallel()

	var _ YouTubeProvider = fakeProvider{}
	var _ TrendsProvider = fakeProvider{}
	var _ LLMProvider = fakeProvider{}
	var _ ImageProvider = fakeProvider{}
	var _ VoiceProvider = fakeProvider{}
	var _ RenderProvider = fakeProvider{}
	var _ UploadProvider = fakeProvider{}
}

type fakeProvider struct{}

func (fakeProvider) SearchVideos(context.Context, YouTubeSearchRequest) (YouTubeSearchResult, error) {
	return YouTubeSearchResult{}, nil
}

func (fakeProvider) GetVideoMetrics(context.Context, VideoMetricsRequest) (VideoMetricsResult, error) {
	return VideoMetricsResult{}, nil
}

func (fakeProvider) InterestOverTime(context.Context, TrendInterestRequest) (TrendInterestResult, error) {
	return TrendInterestResult{}, nil
}

func (fakeProvider) RelatedQueries(context.Context, RelatedQueriesRequest) (RelatedQueriesResult, error) {
	return RelatedQueriesResult{}, nil
}

func (fakeProvider) GenerateText(context.Context, TextGenerationRequest) (TextGenerationResult, error) {
	return TextGenerationResult{}, nil
}

func (fakeProvider) GenerateImage(context.Context, ImageGenerationRequest) (ImageAsset, error) {
	return ImageAsset{}, nil
}

func (fakeProvider) SynthesizeSpeech(context.Context, SpeechSynthesisRequest) (AudioAsset, error) {
	return AudioAsset{}, nil
}

func (fakeProvider) RenderVideo(context.Context, RenderVideoRequest) (VideoAsset, error) {
	return VideoAsset{}, nil
}

func (fakeProvider) UploadPrivateVideo(context.Context, UploadVideoRequest) (UploadVideoResult, error) {
	return UploadVideoResult{}, nil
}
