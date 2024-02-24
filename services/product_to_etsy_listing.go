package services

import (
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
	"github.com/jeremyauchter/uplist/util"
)

type ProductToEtsyListingService struct {
	propertyIds    models.ListingPropertyIDs
	propertyValues models.ListingPropertyValues
	option1Name    string
	option2Name    string
	option3Name    string
	images         models.ListingImages
}

func NewProductToEtsyListingService() *ProductToEtsyListingService {
	propertyValues := make(map[string]int)
	propertyIds := make(map[string]int)
	images := make(map[string]models.ListingImage)
	return &ProductToEtsyListingService{
		propertyIds:    propertyIds,
		propertyValues: propertyValues,
		images:         images,
	}
}

func (s *ProductToEtsyListingService) ConvertToEtsyListing(config *config.Config) (listings []models.EtsyListing, images models.ListingImages, err error) {
	// TODO: Implement the logic to convert a product to an Etsy listing
	productRepo := repositories.NewProductCSVRepository()

	lines := productRepo.GetProducts()
	var listing models.EtsyListing
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
			s.option1Name = line.Option1Name
			s.option2Name = line.Option2Name
			s.option3Name = line.Option3Name

			// build variants
			variant := s.HandleVariant(line, models.EtsyVariant{}, models.EtsyListingInventoryRequest{})
			log.Println("variant", variant)
			listing.Variants = append(listing.Variants, variant)

			// build images
			imageLinks := line.ImageLinks
			mainImage := models.ListingImage{Url: imageLinks, Path: fmt.Sprintf("../../tmp/images/%s", filepath.Base(imageLinks)), Image: filepath.Base(imageLinks)}
			s.images[imageLinks] = mainImage
			listing.Images = append(listing.Images, imageLinks)
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
				s.images[link] = image
				listing.Images = append(listing.Images, link)
			}

		} else {
			variant := s.HandleVariant(line, models.EtsyVariant{}, models.EtsyListingInventoryRequest{})
			log.Println("variant", variant)
			listing.Variants = append(listing.Variants, variant)
		}
	}
	if len(listing.Title) > 0 {
		listings = append(listings, listing)
	}

	log.Println("propertyValues", s.propertyValues)
	// util.PrintJSON(listings)

	return listings, images, nil
}

func (s *ProductToEtsyListingService) HandleVariant(line models.ProductCSV, variant models.EtsyVariant, etsyInventoryRequest models.EtsyListingInventoryRequest) models.EtsyVariant {

	variant.SKU = line.SKU
	variant.Price = line.Price
	variant.Quantity = line.Quantity
	if s.option1Name != "" {
		variant.Option1Name = s.option1Name
		if s.propertyIds[s.option1Name] == 0 {
			s.propertyIds[s.option1Name] = len(s.propertyIds) + 100
		}
		variant.Option1Value = line.Option1Value
		if s.propertyValues[line.Option1Value] == 0 {
			s.propertyValues[line.Option1Value] = len(s.propertyValues) + 1000
		}
	}
	if s.option2Name != "" {
		variant.Option2Name = s.option2Name
		if s.propertyIds[s.option2Name] == 0 {
			s.propertyIds[s.option2Name] = len(s.propertyIds) + 100
		}
		variant.Option2Value = line.Option2Value
		if s.propertyValues[line.Option2Value] == 0 {
			s.propertyValues[line.Option2Value] = len(s.propertyValues) + 1000
		}

	}
	if s.option3Name != "" {
		variant.Option3Name = s.option3Name
		if s.propertyIds[s.option3Name] == 0 {
			s.propertyIds[s.option3Name] = len(s.propertyIds) + 100
		}
		variant.Option3Value = line.Option3Value
		if s.propertyValues[line.Option3Value] == 0 {
			s.propertyValues[line.Option3Value] = len(s.propertyValues) + 1000
		}

	}
	log.Println(s.propertyIds)
	return variant
}

func (s *ProductToEtsyListingService) ConvertVariantsToEtsyProduct(variants []models.EtsyVariant) (etsyInventoryRequest models.EtsyListingInventoryRequest) {
	for _, variant := range variants {
		var product models.EtsyProduct
		product.SKU = variant.SKU
		var productOffering models.EtsyOffering
		productOffering.IsEnabled = true
		productOffering.Price = variant.Price
		productOffering.Quantity = variant.Quantity
		product.Offerings = append(product.Offerings, productOffering)
		if variant.Option1Name != "" {
			var propertyValue models.EtsyPropertyValue
			propertyValue.PropertyID = s.propertyIds[variant.Option1Name]
			propertyValue.PropertyName = variant.Option1Name
			propertyValue.ValueIDs = append(propertyValue.ValueIDs, s.propertyValues[variant.Option1Value])
			propertyValue.Values = append(propertyValue.Values, variant.Option1Value)
			product.PropertyValues = append(product.PropertyValues, propertyValue)
		}
		if variant.Option2Name != "" {
			var propertyValue models.EtsyPropertyValue
			propertyValue.PropertyID = s.propertyIds[variant.Option2Name]
			propertyValue.PropertyName = variant.Option2Name
			propertyValue.ValueIDs = append(propertyValue.ValueIDs, s.propertyValues[variant.Option2Value])
			propertyValue.Values = append(propertyValue.Values, variant.Option2Value)
			product.PropertyValues = append(product.PropertyValues, propertyValue)
		}
		if variant.Option3Name != "" {
			var propertyValue models.EtsyPropertyValue
			propertyValue.PropertyID = s.propertyIds[variant.Option3Name]
			propertyValue.PropertyName = variant.Option3Name
			propertyValue.ValueIDs = append(propertyValue.ValueIDs, s.propertyValues[variant.Option3Value])
			propertyValue.Values = append(propertyValue.Values, variant.Option3Value)
			product.PropertyValues = append(product.PropertyValues, propertyValue)
		}
		etsyInventoryRequest.Products = append(etsyInventoryRequest.Products, product)
		if !util.ContainsInt(etsyInventoryRequest.PriceOnProperty, s.propertyIds[variant.Option1Name]) {
			etsyInventoryRequest.PriceOnProperty = append(etsyInventoryRequest.PriceOnProperty, s.propertyIds[variant.Option1Name])
		}
		if !util.ContainsInt(etsyInventoryRequest.SkuOnProperty, s.propertyIds[variant.Option1Name]) {
			etsyInventoryRequest.SkuOnProperty = append(etsyInventoryRequest.SkuOnProperty, s.propertyIds[variant.Option1Name])
		}
	}

	return etsyInventoryRequest
}

func (s *ProductToEtsyListingService) ConvertImagesToEtsyImageRequests(images []string) (etsyImages []models.EtsyListingImageRequest) {
	for _, image := range images {
		etsyImage := models.EtsyListingImageRequest{
			Image: s.images[image].Path,
		}
		etsyImages = append(etsyImages, etsyImage)
	}
	return etsyImages
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
