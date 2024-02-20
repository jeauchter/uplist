package main

import (
	"fmt"
	"log"

	"github.com/jeremyauchter/uplist/pkg/client"
	"github.com/jeremyauchter/uplist/pkg/config"
	"github.com/jeremyauchter/uplist/services"
)

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
	listings, images, err := etsyProductService.ConvertToEtsyListing(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Println(listings)

	// download all images
	etsyProductService.DownloadImages(etsyAPI, images) // Discard the return value

	// for each product in results
	// Submit Listing to Etsy
	// Submit Images for the Listing
	// Submit Inventory for the Listing
	// update listing to active
	// delete images from local
	etsyProductService.DeleteLocalImages(images)

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
