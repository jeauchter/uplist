package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	APIURL            string `json:"api_url"`
	APIKey            string `json:"api_key"`
	ShopID            string `json:"shop_id"`
	APISecret         string `json:"api_secret"`
	Scopes            string `json:"scopes"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	ShippingProfileID string `json:"shipping_profile_id"`

	// Add more fields based on your configuration file
}

func NewConfig() *Config {
	// Print the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Current working directory:", dir)

	// Continue with your code...

	data, err := os.ReadFile("../../config.json")
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
