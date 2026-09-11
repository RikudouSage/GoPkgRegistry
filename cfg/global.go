package cfg

import "github.com/kelseyhightower/envconfig"

type GlobalConfig struct {
	AdminAPIKey  string `required:"true" split_words:"true"`
	Port         uint16 `default:"8080" split_words:"true"`
	HostOverride string `split_words:"true"`
	DatabasePath string `default:"./data.sqlite3" split_words:"true"`
}

func GetGlobalConfig() (*GlobalConfig, error) {
	cfg := GlobalConfig{}
	err := envconfig.Process("app", &cfg)

	return &cfg, err
}
