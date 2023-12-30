package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jeremyauchter/uplist/repository"
)

type Config struct {
	APIURL string `json:"api_url"`
	// Add more fields based on your configuration file
}

func main() {
	// Read the configuration file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the configuration file
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	resources := repository.GetResources()
	jsonData, err := json.Marshal(resources)
	log.Println(string(jsonData))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	resp, err := http.Post(config.APIURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}
