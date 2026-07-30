package secrets

import (
	"encoding/json"
	"fmt"
	"os"
)

// DefaultPath is where the PaaS automounts the service's secrets in JSON form.
const DefaultPath = "/secrets/secrets.json"

// Secrets holds every secret the Discord bot needs at runtime. It is loaded from
// a single JSON file so nothing sensitive has to be baked into the image.
//
// Note: the YouTube live mode still reads its OAuth credentials from token.json
// directly and is not run from the Docker image, so it is intentionally absent
// here.
type Secrets struct {
	// DiscordToken is the bot token used to authenticate with Discord.
	DiscordToken string `json:"discordToken"`
	// WeatherApiToken is the OpenWeatherMap API key used by the !weather command.
	WeatherApiToken string `json:"weatherApiToken"`
}

// Load reads and parses the secrets file. The path defaults to DefaultPath but
// can be overridden with the SECRETS_PATH environment variable, which is handy
// for local development.
func Load() (*Secrets, error) {
	path := os.Getenv("SECRETS_PATH")
	if path == "" {
		path = DefaultPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading secrets file %q: %w", path, err)
	}

	var s Secrets
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing secrets file %q: %w", path, err)
	}

	return &s, nil
}
