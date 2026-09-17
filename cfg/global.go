package cfg

import (
	"slices"

	"github.com/kelseyhightower/envconfig"
)

type RedirectTo string

const (
	RedirectToSource RedirectTo = "source"
	RedirectToGo     RedirectTo = "go"
)

func (receiver RedirectTo) GetValidValues() []RedirectTo {
	return []RedirectTo{RedirectToSource, RedirectToGo}
}

func (receiver RedirectTo) IsValid() bool {
	return slices.Contains(receiver.GetValidValues(), receiver)
}

func (receiver RedirectTo) GetDefault() RedirectTo {
	return RedirectToSource
}

type GlobalConfig struct {
	// AdminAPIKey is the secret required to access administrative endpoints.
	AdminAPIKey string `required:"true" split_words:"true"`
	// Port is the TCP port on which the HTTP server listens.
	Port uint16 `default:"8080" split_words:"true"`
	// HostOverride replaces the request host when constructing package metadata.
	HostOverride string `split_words:"true"`
	// DatabasePath is the path to the SQLite database file.
	DatabasePath string `default:"./data.sqlite3" split_words:"true"`
	// FrontendURL configures which origin can access this backend using CORS
	FrontendURL string `split_words:"true"`
	// RedirectTo chooses where does the pkg url redirect to, either "source" (to the source repo) or "go" (to pkg.go.dev)
	RedirectTo RedirectTo `split_words:"true" default:"source"`
}

// GetGlobalConfig loads GlobalConfig from APP_-prefixed environment variables.
// Default values are applied to unset optional fields.
func GetGlobalConfig() (*GlobalConfig, error) {
	cfg := GlobalConfig{}
	err := envconfig.Process("app", &cfg)

	return &cfg, err
}
