package config

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadFromEnvUsesDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := LoadFromEnv(mapLookup(nil), Options{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Environment != "development" {
		t.Fatalf("expected development environment, got %q", cfg.Environment)
	}

	if cfg.Server.Addr != defaultHTTPAddr {
		t.Fatalf("expected default HTTP addr %q, got %q", defaultHTTPAddr, cfg.Server.Addr)
	}
}

func TestLoadFromEnvReadsConfiguredValues(t *testing.T) {
	t.Parallel()

	cfg, err := LoadFromEnv(mapLookup(map[string]string{
		"AUTOTUBE_ENV":       "test",
		"AUTOTUBE_HTTP_ADDR": "127.0.0.1:9090",
		"DATABASE_URL":       "postgres://autotube:autotube@localhost:5432/autotube?sslmode=disable",
		"YOUTUBE_API_KEY":    "youtube",
	}), Options{})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Environment != "test" {
		t.Fatalf("expected test environment, got %q", cfg.Environment)
	}

	if cfg.Server.Addr != "127.0.0.1:9090" {
		t.Fatalf("expected configured HTTP addr, got %q", cfg.Server.Addr)
	}

	if cfg.Database.URL == "" {
		t.Fatal("expected database URL to be read")
	}

	if cfg.Providers.YouTubeAPIKey != "youtube" {
		t.Fatalf("expected provider secret to be read, got %q", cfg.Providers.YouTubeAPIKey)
	}
}

func TestLoadFromEnvRequiresDatabaseWhenRequested(t *testing.T) {
	t.Parallel()

	_, err := LoadFromEnv(mapLookup(nil), Options{RequireDatabase: true})
	if err == nil {
		t.Fatal("expected error")
	}

	var missing MissingEnvError
	if !errors.As(err, &missing) {
		t.Fatalf("expected MissingEnvError, got %T", err)
	}

	if got := err.Error(); got != "database setup error: missing required environment variables: DATABASE_URL" {
		t.Fatalf("unexpected error: %q", got)
	}
}

func TestLoadFromEnvRequiresProviderSecretsWhenRequested(t *testing.T) {
	t.Parallel()

	_, err := LoadFromEnv(mapLookup(map[string]string{
		"YOUTUBE_API_KEY": "youtube",
	}), Options{RequireProviderSecrets: true})
	if err == nil {
		t.Fatal("expected error")
	}

	for _, key := range []string{
		"IMAGE_API_KEY",
		"LLM_API_KEY",
		"RENDER_PROVIDER_KEY",
		"TRENDS_API_KEY",
		"UPLOAD_PROVIDER_KEY",
		"VOICE_API_KEY",
	} {
		if !strings.Contains(err.Error(), key) {
			t.Fatalf("expected error to mention %s, got %q", key, err.Error())
		}
	}
}

func TestLoadFromEnvRejectsInvalidDatabaseURL(t *testing.T) {
	t.Parallel()

	_, err := LoadFromEnv(mapLookup(map[string]string{
		"DATABASE_URL": "://not-a-url",
	}), Options{})
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "invalid DATABASE_URL") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func mapLookup(values map[string]string) LookupEnv {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
