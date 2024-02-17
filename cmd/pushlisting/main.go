package main

import (
	"fmt"
	"log"

	"github.com/jeremyauchter/uplist/pkg/client"
	"github.com/jeremyauchter/uplist/pkg/config"
	"github.com/jeremyauchter/uplist/services"
)

type Config struct {
	APIURL string `json:"api_url"`
	// Add more fields based on your configuration file
}

func main() {
	// Read the configuration file
	config := config.NewConfig()
	// Ping test to the Etsy Open API
	etsyAPI := client.NewEtsyAPI(*config)
	reply, err := etsyAPI.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// Fetch resources from the csv file
	etsyProductService := services.NewProductToEtsyListingService()
	err = etsyProductService.ConvertToEtsyListing(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// log.Println(resources)
	// jsonData, err := json.Marshal(resources)
	// log.Println(string(jsonData))
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// resp, err := http.Post(config.APIURL, "application/json", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	fmt.Println("response :", reply)
}
