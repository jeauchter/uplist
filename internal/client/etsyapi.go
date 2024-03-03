package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/jeremyauchter/uplist/models"
)

type EtsyAPI struct {
	oauth2            *oauth2.Config
	apiKey            string
	apiSecret         string
	shopID            string
	accessToken       string
	refreshToken      string
	expiresAt         time.Time
	shippingProfileID int
}

func NewEtsyAPI(apiKey string, apiSecret string) *EtsyAPI {
	return &EtsyAPI{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (api *EtsyAPI) SetOAuthConfig() *oauth2.Config {
	oauth2 := &oauth2.Config{
		ClientID:     api.apiKey,
		ClientSecret: api.apiSecret,
		Scopes:       []string{strings.Join(strings.Split("listings_r,listings_w,shops_r,shops_w,listings_d", ","), " ")},
		Endpoint: oauth2.Endpoint{
			AuthURL:  " https://www.etsy.com/oauth/connect",
			TokenURL: "https://api.etsy.com/v3/public/oauth/token",
		},
		RedirectURL: "https://d958b797cd46b7a14b1667b41b369d7e.m.pipedream.net",
	}
	api.oauth2 = oauth2
	return api.oauth2
}

func (api *EtsyAPI) SetShopId(shopId string) {
	api.shopID = shopId
}

func (api *EtsyAPI) SetAccessToken(accessToken string) {
	api.accessToken = accessToken
}

func (api *EtsyAPI) SetRefreshToken(refreshToken string) {
	api.refreshToken = refreshToken
}

func (api *EtsyAPI) SetExpiresAt(expiresAt time.Time) {
	api.expiresAt = expiresAt
}

func (api *EtsyAPI) Ping() (int, error) {
	// TODO: Implement the logic to ping the Etsy Open API
	// You can use the accessToken field to authenticate the request

	// Example code to make an HTTP request
	url := "https://api.etsy.com/v3/application/openapi-ping"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse body of response: %v", err)
	}

	return resp.StatusCode, nil
}

func (api *EtsyAPI) GetListingsByShop() (models.EtsyListings, error) {
	baseUrl := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings", api.shopID)

	values := url.Values{}
	values.Add("limit", "10")

	urlWithParams := baseUrl + "?" + values.Encode()

	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return models.EtsyListings{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyListings{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EtsyListings{}, fmt.Errorf("failed to parse body of response: %v", err)
	}

	var listings models.EtsyListings
	err = json.Unmarshal(body, &listings)
	if err != nil {
		return models.EtsyListings{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return listings, nil
}

func (api *EtsyAPI) GetReturnPolicies() (models.EtsyReturnPolicyRepsonse, error) {
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/policies/return", api.shopID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.EtsyReturnPolicyRepsonse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyReturnPolicyRepsonse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EtsyReturnPolicyRepsonse{}, fmt.Errorf("failed to parse body of response: %v", err)
	}

	var returnPolicy models.EtsyReturnPolicyRepsonse
	err = json.Unmarshal(body, &returnPolicy)
	if err != nil {
		return models.EtsyReturnPolicyRepsonse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return returnPolicy, nil
}

func (api *EtsyAPI) GetShippingProfiles() (models.EtsyShippingProfileResponse, error) {
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/shipping-profiles", api.shopID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.EtsyShippingProfileResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyShippingProfileResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.EtsyShippingProfileResponse{}, fmt.Errorf("failed to parse body of response: %v", err)
	}

	var shippingProfiles models.EtsyShippingProfileResponse
	err = json.Unmarshal(body, &shippingProfiles)
	if err != nil {
		return models.EtsyShippingProfileResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return shippingProfiles, nil
}

func (api *EtsyAPI) SubmitListing(listingData models.EtsyListingRequest) (listingResponse models.EtsyListingResponse, err error) {
	// TODO: Implement the logic to submit the listing to the Etsy Open API
	// You can use the accessToken field to authenticate the request
	// and the listingData parameter to provide the listing details

	// Example code to make an HTTP request
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings", api.shopID)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	// TODO: Set the necessary request body and headers for listing submission

	jsonListing, err := json.Marshal(listingData)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Body = io.NopCloser(strings.NewReader(string(jsonListing)))

	req.Header.Set("Content-Type", "application/json")

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
	if resp.StatusCode >= 400 {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to send request: %s", string(body))
	}

	err = json.Unmarshal(body, &listingResponse)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// TODO: Handle the response from the Etsy Open API

	return listingResponse, nil
}

func (api *EtsyAPI) UploadImage(imageData models.EtsyListingImageRequest, listingID int, counter int) (imageReponse models.EtsyListingImageResponse, err error) {

	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings/%s/images", api.shopID, strconv.Itoa(listingID))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Add rank to the multipart form data
	rankPart, err := writer.CreateFormField("rank")
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create form field: %v", err)
	}
	_, err = rankPart.Write([]byte(strconv.Itoa(counter)))
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to write to form field: %v", err)
	}

	// Add image files to the multipart form data
	part, err := writer.CreateFormFile("image", imageData.Image)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create form file: %v", err)
	}
	file, err := os.Open(imageData.Image)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to copy file: %v", err)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Set the request body as the multipart form data
	req.Body = io.NopCloser(body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseBody := bytes.NewBuffer(bodyBytes)

	err = json.Unmarshal(responseBody.Bytes(), &imageReponse)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return imageReponse, nil
}

func (api *EtsyAPI) SubmitInventory(inventoryData models.EtsyListingInventoryRequest, listingID int) (err error) {
	url := fmt.Sprintf("https://api.etsy.com/v3/application/listings/%d/inventory", listingID)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)

	jsonListing, err := json.Marshal(inventoryData)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Body = io.NopCloser(strings.NewReader(string(jsonListing)))

	req.Header.Set("Content-Type", "application/json")

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send request: %s", string(body))
	}

	return nil
}

func (api *EtsyAPI) UpdateListingState(listingID int, state string) (err error) {
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings/%d", api.shopID, listingID)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	req.Header.Add("Authorization", "Bearer "+api.accessToken)

	jsonListing, err := json.Marshal(models.EtsyListingStateRequest{State: state})
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Body = io.NopCloser(strings.NewReader(string(jsonListing)))

	req.Header.Set("Content-Type", "application/json")

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send request: %s", string(body))
	}

	return nil
}

func (api *EtsyAPI) AuthorizeApp() (string, string, time.Time, error) {
	fmt.Println("Authorizing app")
	ctx := context.Background()
	conf := api.oauth2
	// Generate a random code verifier
	codeVerifierBytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, codeVerifierBytes)
	if err != nil {
		log.Fatal(err)
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(codeVerifierBytes)

	// Create the code challenge by hashing the code verifier
	codeChallengeBytes := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(codeChallengeBytes[:])

	// Create the authorization URL with the code challenge
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	fmt.Println("After Granting access to Uplist, you will be redirected to the redirect URL.")
	fmt.Println("Enter the code in the URL: ")
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Println(code)
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", string(codeVerifier)))
	if err != nil {
		log.Println(code, codeVerifier)
		log.Fatal(err)
	}
	// log.Println(tok.AccessToken)

	// log.Println(tok.RefreshToken)
	// log.Println(tok.Expiry)

	client := conf.Client(ctx, tok)
	client.Get("...")
	return tok.AccessToken, tok.RefreshToken, tok.Expiry, nil
}

func (api *EtsyAPI) RefreshToken() (string, string, time.Time, error) {
	ctx := context.Background()
	conf := api.oauth2
	tok := &oauth2.Token{
		AccessToken:  api.accessToken,
		RefreshToken: api.refreshToken,
		TokenType:    "Bearer",
		Expiry:       api.expiresAt,
	}
	// Create a new context with the token source
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{})
	// cts := &customTokenSource{ctx, conf, tok}

	//newTok, err := cts.EtsyToken()
	newTok, err := conf.TokenSource(ctx, tok).Token()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("New access token: %s\n", newTok.AccessToken)
	// fmt.Printf("New refresh token: %s\n", newTok.RefreshToken)
	// fmt.Printf("New expiry: %s\n", newTok.Expiry)
	return newTok.AccessToken, newTok.RefreshToken, newTok.Expiry, nil
}

func (api *EtsyAPI) GetShopID(shopName string) (string, error) {
	urlString := "https://api.etsy.com/v3/application/shops"
	values := url.Values{}
	values.Add("shop_name", shopName)

	urlWithParams := urlString + "?" + values.Encode()
	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.apiKey)
	// req.Header.Add("Authorization", "Bearer "+api.accessToken)

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
	return string(body), nil
}

func (api *EtsyAPI) DownloadImage(image models.ListingImage) error {
	// download image based on url and store it locally in the path

	// Example code to make an HTTP request
	url := image.Url
	filepath := image.Path

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a new file to store the image
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
