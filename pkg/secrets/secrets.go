package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultDir is where the PaaS automounts individual secret files.
const DefaultDir = "/secrets"

// Secrets holds every secret the Discord bot needs at runtime. It is loaded from
// individual files in the secret directory so nothing sensitive has to be baked
// into the image.
//
// Note: the YouTube live mode still reads its OAuth credentials from token.json
// directly and is not run from the Docker image, so it is intentionally absent
// here.
type Secrets struct {
	// DiscordToken is the bot token used to authenticate with Discord.
	DiscordToken string
	// WeatherApiToken is the OpenWeatherMap API key used by the !weather command.
	WeatherApiToken string
}

// Load reads the secrets from individual files in the secret directory.
// The directory defaults to DefaultDir (/secrets) but can be overridden with
// the SECRETS_DIR or SECRETS_PATH environment variable for local development.
func Load() (*Secrets, error) {
	dir := os.Getenv("SECRETS_DIR")
	if dir == "" {
		dir = os.Getenv("SECRETS_PATH")
	}
	if dir == "" {
		dir = DefaultDir
	}

	discordToken, err := readSecretFile(dir, "discordToken", "DISCORD_TOKEN", "discord.token")
	if err != nil {
		return nil, fmt.Errorf("loading DiscordToken from %q: %w", dir, err)
	}

	weatherApiToken, err := readSecretFile(dir, "weatherApiToken", "WEATHER_API_TOKEN", "weatherApi.token")
	if err != nil {
		return nil, fmt.Errorf("loading WeatherApiToken from %q: %w", dir, err)
	}

	return &Secrets{
		DiscordToken:    discordToken,
		WeatherApiToken: weatherApiToken,
	}, nil
}

func readSecretFile(dir string, possibleNames ...string) (string, error) {
	var errs []string
	for _, name := range possibleNames {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimSpace(string(data)), nil
		}
		errs = append(errs, err.Error())
	}
	return "", fmt.Errorf("failed to read any of %v (errors: %s)", possibleNames, strings.Join(errs, "; "))
}

