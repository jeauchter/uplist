package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jeremyauchter/uplist/models"
	"github.com/jeremyauchter/uplist/pkg/config"
)

type EtsyAPI struct {
	config *config.Config
}

func NewEtsyAPI(config config.Config) *EtsyAPI {
	return &EtsyAPI{
		config: &config,
	}
}

func (api *EtsyAPI) Ping() (string, error) {
	// TODO: Implement the logic to ping the Etsy Open API
	// You can use the accessToken field to authenticate the request

	// Example code to make an HTTP request
	url := "https://api.etsy.com/v3/application/openapi-ping"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.config.APIKey)

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(body))

	// TODO: Handle the response from the Etsy Open API

	return string(body), nil
}

func (api *EtsyAPI) SubmitListing(listingData models.EtsyListingRequest) (listingResponse models.EtsyListingResponse, err error) {
	// TODO: Implement the logic to submit the listing to the Etsy Open API
	// You can use the accessToken field to authenticate the request
	// and the listingData parameter to provide the listing details

	// Example code to make an HTTP request
	url := "https://api.etsy.com/v2/listings"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.config.AccessToken)

	// TODO: Set the necessary request body and headers for listing submission

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &listingResponse)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// TODO: Handle the response from the Etsy Open API

	return listingResponse, nil
}

func (api *EtsyAPI) UploadImage(imageData models.EtsyListingImageRequest, listingID int) (imageReponse models.EtsyListingImageResponse, err error) {
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings/%s/images", api.config.ShopID, strconv.Itoa(listingID))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.config.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &imageReponse)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return imageReponse, nil
}
