package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	APIURL      string `json:"api_url"`
	APIKey      string `json:"api_key"`
	AccessToken string `json:"access_token"`
	ShopID      string `json:"shop_id"`
	// Add more fields based on your configuration file
}

func NewConfig() *Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the configuration file
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}

func (c *Config) SetAPIURL(url string) {
	c.APIURL = url
}
