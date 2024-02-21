package main

import (
	"fmt"
	"log"

	"github.com/jeremyauchter/uplist/pkg/client"
	"github.com/jeremyauchter/uplist/pkg/config"
	"github.com/jeremyauchter/uplist/services"
	"github.com/jeremyauchter/uplist/util"
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
	listings, _, err := etsyProductService.ConvertToEtsyListing(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	util.PrintJSON(listings)

	// download all images
	// etsyProductService.DownloadImages(etsyAPI, images) // Discard the return value

	// for each product in results
	for _, listing := range listings {
		// Submit Listing to Etsy
		// Submit Images for the Listing
		// Submit Inventory for the Listing
		inventoryRequest := etsyProductService.ConvertVariantsToEtsyProduct(listing.Variants)
		util.PrintJSON(inventoryRequest)
		// util.PrintJSON(listing.Variants)
		// update listing to active
	}
	// delete images from local
	// etsyProductService.DeleteLocalImages(images)

	fmt.Println("response :", reply)
}
