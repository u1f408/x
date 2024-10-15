package main

import (
    "os"
	"context"

    "gopkg.in/yaml.v2"
    "github.com/creasty/defaults"
	"github.com/sethvargo/go-envconfig"
)

func PointerOf[T any](t T) *T {
	return &t
}

type AuthConfig struct {
	Enable bool `yaml:"enable" env:"ENABLE, overwrite" default:"true"`
	HttpHeader string `yaml:"http_header" env:"HTTP_HEADER, overwrite" default:"X-Auth-Request-Email"`
	AllowedDomains []string `yaml:"allowed_domains" env:"ALLOWED_DOMAINS, delimiter=;"`
}

func (c *AuthConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    defaults.Set(c)

    type plain AuthConfig
    if err := unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

type UppyConfig struct {
	ListenHost string `yaml:"listen_host" env:"LISTEN_HOST, overwrite" default:"0.0.0.0"`
	ListenPort uint16 `yaml:"listen_port" env:"LISTEN_PORT, overwrite" default:"23032"`
	Auth AuthConfig `yaml:"authentication" env:", prefix=AUTH_"`
	BunnyCdn BunnyCdnConfig `yaml:"bunnycdn" env:", prefix=BUNNYCDN_"`
}

func (c *UppyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    defaults.Set(c)

    type plain UppyConfig
    if err := unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

func ParseConfigFromFile(cfg *UppyConfig, path string) error {
    raw, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    config := UppyConfig{}
    err = yaml.Unmarshal([]byte(raw), &config)
    if err != nil {
        return err
    }

    *cfg = config
    return nil
}

func ParseConfigFromEnv(cfg *UppyConfig) error {
	config := UppyConfig{}
	if err := defaults.Set(&config); err != nil {
		return err
	}

	if err := envconfig.Process(context.Background(), &config); err != nil {
		return err
	}

	*cfg = config
	return nil
}
