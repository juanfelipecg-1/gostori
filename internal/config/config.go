package config

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	EnvironmentName string `mapstructure:"ENVIRONMENT_NAME"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUsername      string `mapstructure:"DB_USERNAME"`
	DBSchema        string `mapstructure:"DB_SCHEMA"`
	SMTPHost        string `mapstructure:"SMTP_HOST"`
	SMTPPort        string `mapstructure:"SMTP_PORT"`
}

func LoadConfig(paths ...string) (*Config, error) {
	v := viper.New()
	for _, path := range paths {
		v.SetConfigFile(filepath.Clean(path))
		if err := v.MergeInConfig(); err != nil {
			var pathError *os.PathError
			if errors.As(err, &pathError) {
				continue
			}
			return nil, err
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
