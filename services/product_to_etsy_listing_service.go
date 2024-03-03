package services

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	goetsyapi "github.com/jeauchter/go-etsy-api"
	etsyapimodels "github.com/jeauchter/go-etsy-api/models"
	"github.com/jeauchter/uplist/models"
	"github.com/jeauchter/uplist/repositories"
	"github.com/jeauchter/uplist/util"
	"github.com/schollz/progressbar/v3"
)

type ProductToEtsyListingService struct {
	propertyIds    etsyapimodels.ListingPropertyIDs
	propertyValues etsyapimodels.ListingPropertyValues
	option1Name    string
	option2Name    string
	option3Name    string
	images         etsyapimodels.ListingImages
}

func NewProductToEtsyListingService() *ProductToEtsyListingService {
	propertyValues := make(map[string]int)
	propertyIds := make(map[string]int)
	images := make(map[string]etsyapimodels.ListingImage)
	return &ProductToEtsyListingService{
		propertyIds:    propertyIds,
		propertyValues: propertyValues,
		images:         images,
	}
}

func (s *ProductToEtsyListingService) CreateTempImageDirectory(tempDir string) error {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err = os.Mkdir(tempDir, 0755)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Println("Temp directory created: " + tempDir)
	} else {
		fmt.Println("Temp directory already exists: " + tempDir)

	}

	return nil
}

func (s *ProductToEtsyListingService) ConvertToEtsyListing(csvPath string, tempStore string) (listings []etsyapimodels.EtsyListing, images etsyapimodels.ListingImages, err error) {
	// TODO: Implement the logic to convert a product to an Etsy listing
	productRepo := repositories.NewProductCSVRepository()

	uplistImageDir := fmt.Sprintf("%s/uplist-images", tempStore)
	err = s.CreateTempImageDirectory(uplistImageDir)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	lines := productRepo.GetProducts(csvPath)
	var listing etsyapimodels.EtsyListing
	for _, line := range lines {
		if line.Title != "" {

			// build listing
			if len(listing.Title) > 0 {
				listings = append(listings, listing)
			}
			listing = etsyapimodels.EtsyListing{}
			listing.Title = line.Title
			listing.Description = line.Description
			listing.Price = line.Price
			listing.Quantity = line.Quantity
			listing.Tags = line.Tags
			s.option1Name = line.Option1Name
			s.option2Name = line.Option2Name
			s.option3Name = line.Option3Name

			// build variants
			variant := s.HandleVariant(line, etsyapimodels.EtsyVariant{}, etsyapimodels.EtsyListingInventoryRequest{})
			listing.Variants = append(listing.Variants, variant)

			// build images
			imageLinks := line.ImageLinks
			mainImage := etsyapimodels.ListingImage{Url: imageLinks, Path: fmt.Sprintf("%s/%s", uplistImageDir, filepath.Base(imageLinks)), Image: filepath.Base(imageLinks), Rank: 1}
			s.images[imageLinks] = mainImage
			listing.Images = append(listing.Images, imageLinks)
			altImages := line.AdditionalImageLink
			for counter, link := range altImages {
				// breakdown the image link
				rank := counter + 1
				parsedURL, err := url.Parse(link)
				if err != nil {
					log.Println("Failed to parse URL:", err)
					continue
				}

				filename := filepath.Base(parsedURL.Path)
				image := etsyapimodels.ListingImage{Url: link, Path: fmt.Sprintf("%s/%s", uplistImageDir, filename), Image: filename, Rank: rank}
				s.images[link] = image
				listing.Images = append(listing.Images, link)
			}

		} else {
			variant := s.HandleVariant(line, etsyapimodels.EtsyVariant{}, etsyapimodels.EtsyListingInventoryRequest{})
			listing.Variants = append(listing.Variants, variant)
		}
	}
	if len(listing.Title) > 0 {
		listings = append(listings, listing)
	}

	return listings, s.images, nil
}

func (s *ProductToEtsyListingService) HandleVariant(line models.ProductCSV, variant etsyapimodels.EtsyVariant, etsyInventoryRequest etsyapimodels.EtsyListingInventoryRequest) etsyapimodels.EtsyVariant {

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
	return variant
}

func (s *ProductToEtsyListingService) ConvertVariantsToEtsyProduct(variants []etsyapimodels.EtsyVariant) (etsyInventoryRequest etsyapimodels.EtsyListingInventoryRequest) {
	for _, variant := range variants {
		var product etsyapimodels.EtsyProduct
		product.SKU = variant.SKU
		var productOffering etsyapimodels.EtsyOffering
		productOffering.IsEnabled = true
		productOffering.Price = variant.Price
		productOffering.Quantity = variant.Quantity
		product.Offerings = append(product.Offerings, productOffering)
		if variant.Option1Name != "" {
			var propertyValue etsyapimodels.EtsyPropertyValue
			propertyValue.PropertyID = s.propertyIds[variant.Option1Name]
			propertyValue.PropertyName = variant.Option1Name
			propertyValue.ValueIDs = append(propertyValue.ValueIDs, s.propertyValues[variant.Option1Value])
			propertyValue.Values = append(propertyValue.Values, variant.Option1Value)
			product.PropertyValues = append(product.PropertyValues, propertyValue)
		}
		if variant.Option2Name != "" {
			var propertyValue etsyapimodels.EtsyPropertyValue
			propertyValue.PropertyID = s.propertyIds[variant.Option2Name]
			propertyValue.PropertyName = variant.Option2Name
			propertyValue.ValueIDs = append(propertyValue.ValueIDs, s.propertyValues[variant.Option2Value])
			propertyValue.Values = append(propertyValue.Values, variant.Option2Value)
			product.PropertyValues = append(product.PropertyValues, propertyValue)
		}
		if variant.Option3Name != "" {
			var propertyValue etsyapimodels.EtsyPropertyValue
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

func (s *ProductToEtsyListingService) ConvertImagesToEtsyImageRequests(images []string) (etsyImages []etsyapimodels.EtsyListingImageRequest) {
	for _, image := range images {
		etsyImage := etsyapimodels.EtsyListingImageRequest{
			Image: s.images[image].Path,
		}
		etsyImages = append(etsyImages, etsyImage)
	}
	return etsyImages
}

func (s *ProductToEtsyListingService) SubmitListingToEtsy(listing etsyapimodels.EtsyListing, etsyApi *goetsyapi.EtsyAPI, returnPolicyID int, shippingProfileId int) (listingId int, listingTitle string, err error) {

	// build base listing
	baseListing := etsyapimodels.EtsyListingRequest{
		Quantity:          listing.Quantity,
		Title:             listing.Title,
		Description:       listing.Description,
		Price:             listing.Price,
		WhoMade:           etsyapimodels.IDid,
		WhenMade:          etsyapimodels.Year2020_2024,
		TaxonomyID:        399,
		Tags:              strings.Join(listing.Tags, ","),
		ShippingProfileID: shippingProfileId,
		Type:              etsyapimodels.Physical,
		ReturnPolicyID:    returnPolicyID,
	}
	listingResponse, err := etsyApi.SubmitListing(baseListing)
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	// get listing ID
	listingID := listingResponse.ListingID
	return listingID, listingResponse.Title, nil
}

func (s *ProductToEtsyListingService) DownloadImages(etsyApi *goetsyapi.EtsyAPI, images etsyapimodels.ListingImages) {
	bar := progressbar.Default(int64(len(images)))
	for _, image := range images {
		bar.Add(1)
		// only download the image if it doesn't exist
		filepath := image.Path
		if _, err := os.Stat(filepath); err != nil {
			etsyApi.DownloadImage(image)
		}
	}
}

func (s *ProductToEtsyListingService) DeleteLocalImages(images etsyapimodels.ListingImages) {
	for _, image := range images {
		filepath := image.Path
		err := os.Remove(filepath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
