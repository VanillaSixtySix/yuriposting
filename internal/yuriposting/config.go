package yuriposting

import (
	"encoding/json"
	"os"
)

type Config struct {
	MastodonHost        string
	MastodonAccessToken string
	DanbooruUsername    string
	DanbooruAPIKey      string
	Tags                string
	Visibility          string
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
