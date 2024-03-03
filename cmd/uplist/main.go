package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/jeremyauchter/uplist/internal/client"
	"github.com/jeremyauchter/uplist/internal/database"
	"github.com/jeremyauchter/uplist/models"
	"github.com/jeremyauchter/uplist/repositories"
	"github.com/jeremyauchter/uplist/services"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

type Uplist struct {
	etsyAPI         *client.EtsyAPI
	csvPath         string
	PersistentStore string
	TempStore       string
	db              *gorm.DB
}

func main() {
	ul := NewUplist()
	ul.init()
	ul.Run()
}

func NewUplist() *Uplist {
	return &Uplist{}
}

func (ul *Uplist) init() {
	ul.DetermineOS()
	fmt.Printf("Persistent Storage set to %s\n", ul.PersistentStore)
	fmt.Println("Connecting to Database")
	var dbFile string = ul.PersistentStore + "/uplist.db"
	db, err := database.InitDB(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		panic(err)
	} else {
		fmt.Println("Database Connected")
	}
	ul.db = db

	var uplistOauth = repositories.NewUplistOauthRepository(ul.db)
	// Assumes at least one record is found, Run through initial setup to create a record
	shopConfig, err := uplistOauth.Read(1)
	if err != nil {
		// log.Fatalf("Failed to read oauth from database: %v", err)
		fmt.Println("Looks like this is a first time setup, please provide the following information")
		fmt.Println("Please enter the Shop Name")
		var shopName string
		if _, err := fmt.Scan(&shopName); err != nil {
			log.Println(shopName)
			log.Fatal(err)
		}
		fmt.Println("Shop Name Configured as: ", shopName)
		apiKey, apiSecret := ul.SetApiKeys()
		newOauthRecord := &models.UplistOauth{
			ShopName:  shopName,
			APIKey:    apiKey,
			APISecret: apiSecret,
		}
		err = uplistOauth.Create(newOauthRecord)
		if err != nil {
			log.Fatalf("Failed to create shop config record: %v", err)
			panic(err)
		}
	} else {
		fmt.Printf("Shop Name Configured as: %s\n", shopConfig.ShopName)
	}
	ul.etsyAPI = client.NewEtsyAPI(shopConfig.APIKey, shopConfig.APISecret)

	statusCode, err := ul.etsyAPI.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if statusCode > 400 {
		fmt.Println("Looks like the API Key and Secret are not valid, please provide the following information")
		apiKey, apiSecret := ul.SetApiKeys()
		shopConfig.APIKey = apiKey
		shopConfig.APISecret = apiSecret
		err = uplistOauth.Update(shopConfig)
		if err != nil {
			log.Fatalf("Failed to update shop config record: %v", err)
			panic(err)
		}
		ul.etsyAPI = client.NewEtsyAPI(shopConfig.APIKey, shopConfig.APISecret)
	}
	statusCode, err = ul.etsyAPI.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if statusCode > 400 {
		fmt.Println("Looks like the API Key and Secret are not valid, please contact me")
		uplistOauth.Delete(1)
		os.Exit(1)
	}
	fmt.Println("Etsy API Connected!")
	if shopConfig.AccessToken == "" {
		fmt.Println("Looks like this is a first time setup, please complete the following instructions")
		authorizeService := services.NewAuthorizeService(uplistOauth)
		ul.etsyAPI.SetOAuthConfig()
		status, err := authorizeService.AuthorizeApp(ul.etsyAPI)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		fmt.Println(status)
	}
	// Check if the shop ID is configured
	if shopConfig.ShopID == 0 {
		fmt.Println("Shop ID not configured, fetching from Etsy")
		results, err := ul.etsyAPI.GetShopID(shopConfig.ShopName)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		resultsBytes := []byte(results)
		shopResults := models.EtsyShop{}
		err = json.Unmarshal(resultsBytes, &shopResults) // Fix: Pass the address of shopResults
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		shopConfig.ShopID = shopResults.Results[0].ShopID // Fix: Assign the address of the integer value
		err = uplistOauth.Update(shopConfig)
		if err != nil {
			log.Fatalf("Failed to update shop config record: %v", err)
			panic(err)
		}
		fmt.Printf("Shop ID Configured as: %d\n", shopConfig.ShopID)
	}
	ul.etsyAPI.SetShopId(strconv.Itoa(shopConfig.ShopID))

}

func (ul *Uplist) SetApiKeys() (string, string) {
	fmt.Println("Please enter the API Key for Uplist")
	var apiKey string
	if _, err := fmt.Scan(&apiKey); err != nil {
		log.Println(apiKey)
		log.Fatal(err)
	}
	fmt.Println("API Key Configured as: ", apiKey)
	fmt.Println("Please enter the API Secret for Uplist")
	var apiSecret string
	if _, err := fmt.Scan(&apiSecret); err != nil {
		log.Println(apiSecret)
		log.Fatal(err)
	}
	fmt.Println("API Secret Configured")
	return apiKey, apiSecret
}

func (ul *Uplist) Run() {
	var uplistOauth = repositories.NewUplistOauthRepository(ul.db)
	authorizeService := services.NewAuthorizeService(uplistOauth)
	ul.etsyAPI.SetOAuthConfig()
	status, err := authorizeService.RefreshToken(ul.etsyAPI)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println(status)
	// get shop listings
	// listings, err := ul.etsyAPI.GetListingsByShop()
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// // util.PrintJSON(listings)

	ul.GetCSVPath()
	ul.publishListings()
}

func (ul *Uplist) DetermineOS() {
	// Determine the operating system
	ul.TempStore = os.TempDir()
	switch runningOS := runtime.GOOS; runningOS {
	case "windows":
		fmt.Println("Determined Uplist is Running on Windows")
		ul.PersistentStore = os.Getenv("USERPROFILE")
	case "linux":
		fmt.Println("Determined Uplist is Running on Linux")
		ul.PersistentStore = os.Getenv("HOME")
	case "darwin":
		fmt.Println("Determined Uplist is Running on macOS")
		ul.PersistentStore = os.Getenv("HOME")
	default:
		fmt.Printf("Running on an unsupported operating system: %s\n", runningOS)
		os.Exit(1)
	}
}

func (ul *Uplist) GetCSVPath() {
	var csvPath string

	flag.StringVar(&csvPath, "csv", "", "Path to CSV")
	flag.Parse()

	if csvPath == "" {
		fmt.Println("Please provide a path to the CSV file")
		os.Exit(1)
	}

	ul.csvPath = csvPath
}

func (ul *Uplist) publishListings() {

	// Fetch resources from the csv file
	etsyProductService := services.NewProductToEtsyListingService()
	listings, images, err := etsyProductService.ConvertToEtsyListing(ul.csvPath, ul.TempStore)
	fmt.Println("Listings Converted! Count: " + strconv.Itoa(len(listings)))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// util.PrintJSON(listings)
	// Get Return Policy
	returnPolicies, err := ul.etsyAPI.GetReturnPolicies()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	returnPolicy := returnPolicies.Results[0].ReturnPolicyID
	fmt.Println("Return Policy: ", returnPolicy)

	shippingProfiles, err := ul.etsyAPI.GetShippingProfiles()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	shippingProfileId := shippingProfiles.Results[0].ShippingProfileID
	// download all images
	etsyProductService.DownloadImages(ul.etsyAPI, images) // Discard the return value

	// for each product in results
	for _, listing := range listings {
		// Submit Listing to Etsy
		listingId, err := etsyProductService.SubmitListingToEtsy(listing, ul.etsyAPI, returnPolicy, shippingProfileId)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		fmt.Println("Listing ID: ", listingId)

		// Submit Images for the Listing
		fmt.Println("Uploading Images")
		listingImages := etsyProductService.ConvertImagesToEtsyImageRequests(listing.Images)
		bar := progressbar.Default(int64(len(listingImages)))
		for counter, image := range listingImages {

			bar.Add(1)
			_, err := ul.etsyAPI.UploadImage(image, listingId, counter+1)
			if err != nil {
				log.Fatal(err)
				panic(err)
				//something went wrong with image upload, delete the listing
				//TODO: ul.etsyAPI.DeleteListing(listingId)
				// ul.etsyAPI.DeleteListing(listingId)
			}

		}
		fmt.Println("Images Uploaded!")
		fmt.Println("Uploading Inventory")
		// Submit Inventory for the Listing
		inventoryRequest := etsyProductService.ConvertVariantsToEtsyProduct(listing.Variants)
		// util.PrintJSON(inventoryRequest)
		err = ul.etsyAPI.SubmitInventory(inventoryRequest, listingId)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		fmt.Println("Inventory Uploaded!")

		// util.PrintJSON(listing.Variants)
		// update listing to active
		fmt.Println("Updating Listing State")
		err = ul.etsyAPI.UpdateListingState(listingId, "active")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		fmt.Println("Listing State Updated!")
		fmt.Println("Listing ID: ", listingId)
	}
	// delete images from local
	etsyProductService.DeleteLocalImages(images)
	fmt.Println("Listings Published!")
}
