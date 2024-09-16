package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/response"
	"net/http"
)

var (
	client = &http.Client{}
)

func FetchFacebookProfile(token string) (*response.FetchFacebookProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=%s", token), nil)
	if err != nil {
		return nil, errors.New("error creating HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error sending HTTP request")
	}

	defer resp.Body.Close()

	responseStruct := new(response.FetchFacebookProfile)
	if err = json.NewDecoder(resp.Body).Decode(responseStruct); err != nil {
		return nil, errors.New("error decode into struct")
	}

	return responseStruct, nil
}

func FetchGoogleProfile(token string) (*response.FetchGoogleProfile, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token), nil)
	if err != nil {
		return nil, errors.New("error creating HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error sending HTTP request")
	}

	defer resp.Body.Close()

	responseStruct := new(response.FetchGoogleProfile)
	if err = json.NewDecoder(resp.Body).Decode(responseStruct); err != nil {
		return nil, errors.New("error decode into struct")
	}

	return responseStruct, nil
}

func DetermineRedirectURL(res *response.OAuth) string {
	if res.Error {
		return fmt.Sprintf("https://example.com/dashboard?&message=%s&error=%t", res.ErrorMessage, res.Error)
	}

	if res.IsRegistered {
		return fmt.Sprintf("https://example.com/dashboard?access_token=%s&is_registered=%t", res.AccessToken, res.IsRegistered)
	}

	return fmt.Sprintf("https://example.com/dashboard?token=%s&is_registered=%t", res.Token, res.IsRegistered)
}
