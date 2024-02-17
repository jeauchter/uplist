package services

import (
	"log"
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

func (s *ProductToEtsyListingService) ConvertToEtsyListing(config *config.Config) error {
	// TODO: Implement the logic to convert a product to an Etsy listing
	productRepo := repositories.NewProductCSVRepository()
	etsyApi := client.NewEtsyAPI(*config)

	lines := productRepo.GetProducts()
	for _, line := range lines {
		if line.Title != "" {

			intShippingProfileID, err := strconv.Atoi(config.ShippingProfileID)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
			log.Println("intShippeingId", intShippingProfileID)

			// build base listing
			baseListing := models.EtsyListingRequest{
				ShopID:            config.ShopID,
				Quantity:          line.Quantity,
				Title:             line.Title,
				Description:       line.Description,
				Price:             line.Price,
				WhoMade:           models.IDid,
				WhenMade:          models.Year2020_2024,
				TaxonomyID:        1,
				Tags:              strings.Join(line.Tags, ","),
				ShippingProfileID: intShippingProfileID,
				Type:              models.Physical,
			}
			log.Println(baseListing)
			listingResponse, err := etsyApi.SubmitListing(baseListing)
			if err != nil {
				log.Fatal(err)
				panic(err)
			}

			// get listing ID
			listingID := listingResponse.ListingID
			log.Println("Listing Id = " + strconv.Itoa(listingID))

			// get product images
			images := []string{}
			images = append(images, line.ImageLinks)
			images = append(images, line.AdditionalImageLink...)
			log.Println(baseListing)
			// build image request
			for _, image := range images {
				//build image request
				var imageRequest models.EtsyListingImageRequest
				imageRequest.Image = image
				log.Println(imageRequest)
				// imageResponse, err := etsyApi.UploadImage(imageRequest, listingID)
				// if err != nil {
				// 	log.Fatal(err)
				// 	panic(err)
				// }
				// log.Println(imageResponse)

			}
		}
	}

	// jsonData, err := json.MarshalIndent(lines, "", "    ")
	// log.Println(string(jsonData))
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }

	return nil
}
