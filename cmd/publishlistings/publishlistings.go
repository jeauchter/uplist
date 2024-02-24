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
	listings, images, err := etsyProductService.ConvertToEtsyListing(config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// util.PrintJSON(listings)

	// download all images
	etsyProductService.DownloadImages(etsyAPI, images) // Discard the return value

	// for each product in results
	for _, listing := range listings {
		// Submit Listing to Etsy
		listingId, err := etsyProductService.SubmitListingToEtsy(listing, etsyAPI, config)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		// Submit Images for the Listing
		listingImages := etsyProductService.ConvertImagesToEtsyImageRequests(listing.Images)
		for _, image := range listingImages {
			_, err := etsyAPI.UploadImage(image, listingId)
			if err != nil {
				log.Fatal(err)
				panic(err)
				//something went wrong with image upload, delete the listing
				//TODO: etsyAPI.DeleteListing(listingId)
				// etsyAPI.DeleteListing(listingId)
			}

		}
		// Submit Inventory for the Listing
		inventoryRequest := etsyProductService.ConvertVariantsToEtsyProduct(listing.Variants)
		util.PrintJSON(inventoryRequest)
		// TODO: implement token store in database over passing it all over the place
		err = etsyAPI.SubmitInventory(inventoryRequest, listingId, "token")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		// util.PrintJSON(listing.Variants)
		// update listing to active
		// TODO: implement UpdateListingState
		// _, err = etsyAPI.UpdateListingState(listingId, "active")
	}
	// delete images from local
	etsyProductService.DeleteLocalImages(images)

	fmt.Println("response :", reply)
}
