package services

import (
	"fmt"
	"log"

	goetsyapi "github.com/jeauchter/go-etsy-api"
	"github.com/jeauchter/uplist/models"
	"github.com/jeauchter/uplist/repositories"
)

type AuthorizeService interface {
	AuthorizeApp(etsyApi *goetsyapi.EtsyAPI) (string, error)
	RefreshToken(etsyApi *goetsyapi.EtsyAPI) (string, error)
}

type AuthorizeServiceImp struct {
	rep *repositories.UplistOauthRepository
}

func NewAuthorizeService(rep *repositories.UplistOauthRepository) AuthorizeService {
	return &AuthorizeServiceImp{
		rep: rep,
	}
}

func (a *AuthorizeServiceImp) AuthorizeApp(etsyApi *goetsyapi.EtsyAPI) (string, error) {
	accessToken, refreshToken, Expiry, err := etsyApi.AuthorizeApp()
	if err != nil {
		log.Fatal(err)
	}
	updateOauthRecord := &models.UplistOauth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    Expiry,
	}
	a.rep.UpdateByID(uint(1), updateOauthRecord)
	etsyApi.SetAccessToken(accessToken)
	etsyApi.SetRefreshToken(refreshToken)
	etsyApi.SetExpiresAt(Expiry)
	return "Authorized!", nil
}

func (a *AuthorizeServiceImp) RefreshToken(etsyApi *goetsyapi.EtsyAPI) (string, error) {
	existing, err := a.rep.Read(1)
	if err != nil {
		log.Fatal(err)
	}
	etsyApi.SetAccessToken(existing.AccessToken)
	etsyApi.SetRefreshToken(existing.RefreshToken)
	etsyApi.SetExpiresAt(existing.ExpiresAt)
	accessToken, refreshToken, Expiry, err := etsyApi.RefreshToken()
	if err != nil {
		log.Fatal(err)
	}
	updateOauthRecord := &models.UplistOauth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    Expiry,
	}
	a.rep.UpdateByID(uint(1), updateOauthRecord)
	etsyApi.SetAccessToken(accessToken)
	etsyApi.SetRefreshToken(refreshToken)
	etsyApi.SetExpiresAt(Expiry)
	return fmt.Sprintf("Refreshed! New Access Token %s", accessToken), nil
}

func (a *AuthorizeServiceImp) GetAccessToken() (string, error) {
	// TODO: Fetch the access token from the database
	shopConfig, err := a.rep.Read(1)
	if err != nil {
		log.Fatal(err)
	}

	return shopConfig.AccessToken, nil
}
