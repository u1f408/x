package main

import (
    "os"
    "strings"

    "gopkg.in/yaml.v2"
    "github.com/creasty/defaults"
)

type ProxyDnsConfig struct {
    Enable bool `yaml:"enable"`
    Host string `yaml:"host"`
    Port string `yaml:"port"`
}

func (c *ProxyDnsConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var err error

    var rawHostPort string
    if err = unmarshal(&rawHostPort); err == nil {
        if len(rawHostPort) == 0 {
            c.Enable = false
            c.Host = ""
            c.Port = ""

            return nil
        }

        c.Enable = true
        lastColon := strings.LastIndex(rawHostPort, ":")
        if lastColon == -1 {
            c.Host = rawHostPort
            c.Port = "53"
        } else {
            c.Host = rawHostPort[:lastColon]
            c.Port = rawHostPort[lastColon:]
        }

        return nil
    }

    type plain ProxyDnsConfig
    if err = unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

type ProxyLookupConsul struct {
    Enable bool `yaml:"enable" default:"false"`
    ServiceAddr string `yaml:"service"`
    LookupDNS ProxyDnsConfig `yaml:"lookup_dns"`
    ProxyDNS bool `yaml:"proxy_dns" default:"true"`
}

func (c *ProxyLookupConsul) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var err error
    defaults.Set(c)

    var rawString string
    if err = unmarshal(&rawString); err == nil && rawString != "" {
        c.Enable = true
        c.ServiceAddr = rawString

        return nil
    }

    type plain ProxyLookupConsul
    if err = unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

type ProxyLookupConfig struct {
    Consul ProxyLookupConsul `yaml:"consul"`
}

type ProxyConfig struct {
    Domain string `yaml:"domain"`
    StripDomain bool `yaml:"strip_domain" default:"false"`
    ProxyUrl string `yaml:"proxy_url"`
    Lookup ProxyLookupConfig `yaml:"lookup"`
}

func (c *ProxyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    defaults.Set(c)

    type plain ProxyConfig
    if err := unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

type Config struct {
    Proxies []ProxyConfig `yaml:"proxies"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
    defaults.Set(c)

    type plain Config
    if err := unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

func ParseConfig(cfg *Config, path string) error {
    raw, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    config := Config{}
    err = yaml.Unmarshal([]byte(raw), &config)
    if err != nil {
        return err
    }

    *cfg = config
    return nil
}
