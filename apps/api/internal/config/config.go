package config

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

const (
	defaultHTTPAddr = ":8080"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Providers   ProviderConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL string
}

type ProviderConfig struct {
	YouTubeAPIKey     string
	TrendsAPIKey      string
	LLMAPIKey         string
	ImageAPIKey       string
	VoiceAPIKey       string
	RenderProviderKey string
	UploadProviderKey string
}

type Options struct {
	RequireDatabase        bool
	RequireProviderSecrets bool
}

type LookupEnv func(string) (string, bool)

type MissingEnvError struct {
	Scope string
	Keys  []string
}

func (e MissingEnvError) Error() string {
	keys := append([]string(nil), e.Keys...)
	sort.Strings(keys)

	scope := strings.TrimSpace(e.Scope)
	if scope == "" {
		scope = "configuration"
	}

	return fmt.Sprintf("%s setup error: missing required environment variables: %s", scope, strings.Join(keys, ", "))
}

func Load(options Options) (Config, error) {
	return LoadFromEnv(os.LookupEnv, options)
}

func LoadFromEnv(lookup LookupEnv, options Options) (Config, error) {
	if lookup == nil {
		lookup = os.LookupEnv
	}

	cfg := Config{
		Environment: getEnv(lookup, "AUTOTUBE_ENV", "development"),
		Server: ServerConfig{
			Addr: getEnv(lookup, "AUTOTUBE_HTTP_ADDR", defaultHTTPAddr),
		},
		Database: DatabaseConfig{
			URL: getEnv(lookup, "DATABASE_URL", ""),
		},
		Providers: ProviderConfig{
			YouTubeAPIKey:     getEnv(lookup, "YOUTUBE_API_KEY", ""),
			TrendsAPIKey:      getEnv(lookup, "TRENDS_API_KEY", ""),
			LLMAPIKey:         getEnv(lookup, "LLM_API_KEY", ""),
			ImageAPIKey:       getEnv(lookup, "IMAGE_API_KEY", ""),
			VoiceAPIKey:       getEnv(lookup, "VOICE_API_KEY", ""),
			RenderProviderKey: getEnv(lookup, "RENDER_PROVIDER_KEY", ""),
			UploadProviderKey: getEnv(lookup, "UPLOAD_PROVIDER_KEY", ""),
		},
	}

	if err := cfg.Validate(options); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate(options Options) error {
	if strings.TrimSpace(cfg.Environment) == "" {
		return MissingEnvError{Scope: "configuration", Keys: []string{"AUTOTUBE_ENV"}}
	}

	if strings.TrimSpace(cfg.Server.Addr) == "" {
		return MissingEnvError{Scope: "configuration", Keys: []string{"AUTOTUBE_HTTP_ADDR"}}
	}

	if options.RequireDatabase {
		if strings.TrimSpace(cfg.Database.URL) == "" {
			return MissingEnvError{Scope: "database", Keys: []string{"DATABASE_URL"}}
		}
	}

	if strings.TrimSpace(cfg.Database.URL) != "" {
		if _, err := url.ParseRequestURI(cfg.Database.URL); err != nil {
			return fmt.Errorf("database setup error: invalid DATABASE_URL: %w", err)
		}
	}

	if options.RequireProviderSecrets {
		if err := cfg.Providers.ValidateRequired(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg ProviderConfig) ValidateRequired() error {
	missing := make([]string, 0, 7)

	if strings.TrimSpace(cfg.YouTubeAPIKey) == "" {
		missing = append(missing, "YOUTUBE_API_KEY")
	}
	if strings.TrimSpace(cfg.TrendsAPIKey) == "" {
		missing = append(missing, "TRENDS_API_KEY")
	}
	if strings.TrimSpace(cfg.LLMAPIKey) == "" {
		missing = append(missing, "LLM_API_KEY")
	}
	if strings.TrimSpace(cfg.ImageAPIKey) == "" {
		missing = append(missing, "IMAGE_API_KEY")
	}
	if strings.TrimSpace(cfg.VoiceAPIKey) == "" {
		missing = append(missing, "VOICE_API_KEY")
	}
	if strings.TrimSpace(cfg.RenderProviderKey) == "" {
		missing = append(missing, "RENDER_PROVIDER_KEY")
	}
	if strings.TrimSpace(cfg.UploadProviderKey) == "" {
		missing = append(missing, "UPLOAD_PROVIDER_KEY")
	}

	if len(missing) > 0 {
		return MissingEnvError{Scope: "provider", Keys: missing}
	}

	return nil
}

func getEnv(lookup LookupEnv, key string, fallback string) string {
	value, ok := lookup(key)
	if !ok {
		return fallback
	}

	return strings.TrimSpace(value)
}
