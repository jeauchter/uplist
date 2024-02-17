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

	"golang.org/x/oauth2"

	"github.com/jeremyauchter/uplist/models"
	"github.com/jeremyauchter/uplist/pkg/config"
)

type EtsyAPI struct {
	config *config.Config
	oauth2 *oauth2.Config
}

func NewEtsyAPI(config config.Config) *EtsyAPI {
	return &EtsyAPI{
		config: &config,
		oauth2: &oauth2.Config{
			ClientID:     config.APIKey,
			ClientSecret: config.APISecret,
			Scopes:       []string{strings.Join(strings.Split(config.Scopes, ","), " ")},
			Endpoint: oauth2.Endpoint{
				AuthURL:  " https://www.etsy.com/oauth/connect",
				TokenURL: "https://api.etsy.com/v3/public/oauth/token",
			},
			RedirectURL: "https://d958b797cd46b7a14b1667b41b369d7e.m.pipedream.net",
		},
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

	accessToken, err := api.RefreshToken()
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to refresh token: %v", err)
	}

	// Example code to make an HTTP request
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings", api.config.ShopID)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("x-api-key", api.config.APIKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	log.Println(req)
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
	log.Println(string(jsonListing))

	req.Header.Set("Content-Type", "application/json")

	// Example code to send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	log.Print(resp)
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to send request: %s", string(body))
	}
	log.Println(string(body))

	err = json.Unmarshal(body, &listingResponse)
	if err != nil {
		return models.EtsyListingResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// TODO: Handle the response from the Etsy Open API

	return listingResponse, nil
}

func (api *EtsyAPI) UploadImage(imageData models.EtsyListingImageRequest, listingID int) (imageReponse models.EtsyListingImageResponse, err error) {

	accessToken, err := api.RefreshToken()
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to refresh token: %v", err)
	}
	url := fmt.Sprintf("https://api.etsy.com/v3/application/shops/%s/listings/%s/images", api.config.ShopID, strconv.Itoa(listingID))

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.config.APIKey)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image files to the multipart form data
	for _, imageFile := range imageData.Image {
		file, err := os.Open(string(imageFile))
		if err != nil {
			return models.EtsyListingImageResponse{}, fmt.Errorf("failed to open image file: %v", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("images[]", string(imageFile))
		if err != nil {
			return models.EtsyListingImageResponse{}, fmt.Errorf("failed to create form file: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return models.EtsyListingImageResponse{}, fmt.Errorf("failed to copy file data: %v", err)
		}
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

func (api *EtsyAPI) AuthorizeApp() (string, error) {
	log.Println("Authorizing app")
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     api.config.APIKey,
		ClientSecret: api.config.APISecret,
		Scopes:       []string{strings.Join(strings.Split(api.config.Scopes, ","), " ")},
		Endpoint: oauth2.Endpoint{
			AuthURL:  " https://www.etsy.com/oauth/connect",
			TokenURL: "https://api.etsy.com/v3/public/oauth/token",
		},
		RedirectURL: "https://d958b797cd46b7a14b1667b41b369d7e.m.pipedream.net",
	}
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
	log.Println(tok.AccessToken)

	log.Println(tok.RefreshToken)

	client := conf.Client(ctx, tok)
	client.Get("...")
	return "ok", nil
}

func (api *EtsyAPI) RefreshToken() (string, error) {
	log.Println("Refreshing token")
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     api.config.APIKey,
		ClientSecret: api.config.APISecret,
		Scopes:       []string{strings.Join(strings.Split(api.config.Scopes, ","), " ")},
		Endpoint: oauth2.Endpoint{
			AuthURL:  " https://www.etsy.com/oauth/connect",
			TokenURL: "https://api.etsy.com/v3/public/oauth/token",
		},
		RedirectURL: "https://d958b797cd46b7a14b1667b41b369d7e.m.pipedream.net",
	}
	tok := &oauth2.Token{
		AccessToken:  api.config.AccessToken,
		RefreshToken: api.config.RefreshToken,
		TokenType:    "Bearer",
	}
	// Create a new context with the token source
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{})
	cts := &customTokenSource{ctx, conf, tok}

	newTok, err := cts.EtsyToken()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New access token: %s\n", newTok.AccessToken)
	return newTok.AccessToken, nil
}

type customTokenSource struct {
	ctx  context.Context
	conf *oauth2.Config
	tok  *oauth2.Token
}

func (c *customTokenSource) EtsyToken() (*oauth2.Token, error) {
	if c.tok.Valid() {
		return c.tok, nil
	}

	v := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.tok.RefreshToken},
		"client_id":     {c.conf.ClientID},
	}

	req, err := http.NewRequest("POST", c.conf.Endpoint.TokenURL, strings.NewReader(v.Encode()))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	c.tok = &token
	return c.tok, nil
}

func (api *EtsyAPI) GetShopID() (string, error) {
	url := "https://api.etsy.com/v3/application/shops"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("x-api-key", api.config.APIKey)
	req.Header.Add("Authorization", "Bearer "+api.config.AccessToken)

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
	return string(body), nil
}
