package yuriposting

import (
	"encoding/json"
	"os"
)

type Config struct {
	PostToBluesky          bool
	PostToMastodon         bool
	BlueskyIdentifier      string
	BlueskyAppPassword     string
	MastodonHost           string
	MastodonAccessToken    string
	MastodonPostVisibility string
	DanbooruUsername       string
	DanbooruAPIKey         string
	DanbooruTags           string
}

func LoadConfig(filename string) (*Config, error) {
	var config *Config
	configBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return config, nil
}
