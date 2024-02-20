package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jeremyauchter/uplist/models"
	"github.com/jeremyauchter/uplist/pkg/client"
	"github.com/jeremyauchter/uplist/pkg/config"
	"github.com/jeremyauchter/uplist/repositories"
)

type ProductToEtsyListingService struct {
}

func NewProductToEtsyListingService() *ProductToEtsyListingService {
	return &ProductToEtsyListingService{}
}

func (s *ProductToEtsyListingService) ConvertToEtsyListing(config *config.Config) (listings []models.EtsyListing, images models.ListingImages, err error) {
	// TODO: Implement the logic to convert a product to an Etsy listing
	productRepo := repositories.NewProductCSVRepository()

	lines := productRepo.GetProducts()
	var option_1_name, option_2_name, option_3_name string
	var listing models.EtsyListing
	images = make(map[string]models.ListingImage)
	for _, line := range lines {
		if line.Title != "" {

			// build listing
			if len(listing.Title) > 0 {
				listings = append(listings, listing)
			}
			listing = models.EtsyListing{}
			listing.Title = line.Title
			listing.Description = line.Description
			listing.Price = line.Price
			listing.Quantity = line.Quantity
			listing.Tags = line.Tags

			option_1_name = line.Option1Name
			option_2_name = line.Option2Name
			option_3_name = line.Option3Name
			variant := models.EtsyVariant{
				SKU:      line.SKU,
				Price:    line.Price,
				Quantity: line.Quantity,
			}
			if option_1_name != "" {
				variant.Option1Name = option_1_name
				variant.Option1Value = line.Option1Value
			}
			if option_2_name != "" {
				variant.Option2Name = option_2_name
				variant.Option2Value = line.Option2Value
			}
			if option_3_name != "" {
				variant.Option3Name = option_3_name
				variant.Option3Value = line.Option3Value
			}
			listing.Variants = append(listing.Variants, variant)

			// build images
			imageLinks := line.ImageLinks
			mainImage := models.ListingImage{Url: imageLinks, Path: fmt.Sprintf("../../tmp/images/%s", filepath.Base(imageLinks)), Image: filepath.Base(imageLinks)}
			images[imageLinks] = mainImage
			altImages := line.AdditionalImageLink
			for _, link := range altImages {
				// breakdown the image link

				parsedURL, err := url.Parse(link)
				if err != nil {
					log.Println("Failed to parse URL:", err)
					continue
				}

				filename := filepath.Base(parsedURL.Path)
				image := models.ListingImage{Url: link, Path: fmt.Sprintf("../../tmp/images/%s", filename), Image: filename}
				images[link] = image
			}

		} else {
			variant := models.EtsyVariant{
				SKU:      line.SKU,
				Price:    line.Price,
				Quantity: line.Quantity,
			}
			if option_1_name != "" {
				variant.Option1Name = option_1_name
				variant.Option1Value = line.Option1Value
			}
			if option_2_name != "" {
				variant.Option2Name = option_2_name
				variant.Option2Value = line.Option2Value
			}
			if option_3_name != "" {
				variant.Option3Name = option_3_name
				variant.Option3Value = line.Option3Value
			}
			listing.Variants = append(listing.Variants, variant)
		}
	}
	if len(listing.Title) > 0 {
		listings = append(listings, listing)
	}

	jsonData, err := json.MarshalIndent(images, "", "    ")
	log.Println(string(jsonData))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return listings, images, nil
}

func (s *ProductToEtsyListingService) SubmitListingToEtsy(listing models.EtsyListing, etsyApi *client.EtsyAPI, config *config.Config) (listingId int, err error) {
	intShippingProfileID, err := strconv.Atoi(config.ShippingProfileID)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	log.Println("intShippeingId", intShippingProfileID)

	// build base listing
	baseListing := models.EtsyListingRequest{
		ShopID:            config.ShopID,
		Quantity:          listing.Quantity,
		Title:             listing.Title,
		Description:       listing.Description,
		Price:             listing.Price,
		WhoMade:           models.IDid,
		WhenMade:          models.Year2020_2024,
		TaxonomyID:        399,
		Tags:              strings.Join(listing.Tags, ","),
		ShippingProfileID: intShippingProfileID,
		Type:              models.Physical,
	}
	log.Println(baseListing)
	listingResponse, err := etsyApi.SubmitListing(baseListing)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	// get listing ID
	listingID := listingResponse.ListingID
	return listingID, nil
}

func (s *ProductToEtsyListingService) DownloadImages(etsyApi *client.EtsyAPI, images models.ListingImages) {
	for _, image := range images {
		log.Println(image.Url)
		// only download the image if it doesn't exist
		filepath := image.Path
		if _, err := os.Stat(filepath); err != nil {
			etsyApi.DownloadImage(image)
		} else {
			log.Println("Image already exists:", filepath)
		}
	}
}

func (s *ProductToEtsyListingService) DeleteLocalImages(images models.ListingImages) {
	for _, image := range images {
		filepath := image.Path
		err := os.Remove(filepath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
