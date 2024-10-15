package main

import (
	"fmt"
	"bytes"
	"net/http"
	"github.com/creasty/defaults"
)

type BunnyCdnConfig struct {
	Enable bool `yaml:"enable" env:"ENABLE, overwrite" default:"false"`
	FinalUrl string `yaml:"final_url" env:"FINAL_URL"`
	StorageEndpoint string `yaml:"storage_endpoint" env:"STORAGE_ENDPOINT, overwrite" default:"https://storage.bunnycdn.com"`
	StorageZone string `yaml:"storage_zone_name" env:"STORAGE_ZONE"`
	AccessKey string `yaml:"access_key" env:"ACCESS_KEY"`
}

func (c *BunnyCdnConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    defaults.Set(c)

    type plain BunnyCdnConfig
    if err := unmarshal((*plain)(c)); err != nil {
        return err
    }

    return nil
}

func (c *BunnyCdnConfig) IsValid() bool {
	return c.Enable &&
		len(c.FinalUrl) > 0 &&
		len(c.StorageEndpoint) > 0 &&
		len(c.StorageZone) > 0 &&
		len(c.AccessKey) > 0
}

func (c *BunnyCdnConfig) UploadFile(content *bytes.Buffer, mimeType, dir, filename string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", c.StorageEndpoint, c.StorageZone, dir, filename)
	req, err := http.NewRequest(http.MethodPut, url, content)
	req.Header.Set("AccessKey", c.AccessKey)
	req.Header.Set("Content-Type", mimeType)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("got HTTP %d from bunny.net API, expected 201", resp.StatusCode)
	}

	return fmt.Sprintf("%s/%s/%s", c.FinalUrl, dir, filename), nil
}
