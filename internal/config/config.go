package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type Postgres struct {
	Host     string
	User     string
	Password string
}

type Otel struct {
	Enable            bool
	CollectorEndpoint string
	ServiceName       string
	InsecureMode      bool
}

type Config struct {
	Postgres *Postgres
	Otel     *Otel
}

func Read() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.SetConfigType("yml")

	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
