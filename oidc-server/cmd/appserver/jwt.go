package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

// gets token from tokenUrl validating token with jwksUrl and returning token & claims
func getTokenFromCode(tokenUrl, jwksUrl, redirectUri, clientID, clientSecret, code string) (*jwt.Token, *jwt.StandardClaims, error) {
	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("client_id", clientID)
	values.Add("client_secret", clientSecret)
	values.Add("redirect_uri", redirectUri)
	values.Add("code", code)

	res, err := http.PostForm(tokenUrl, values)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code was not 200")
	}
	fmt.Printf("body: %s", body)
	return nil, nil, fmt.Errorf("not yet fully implemented")
}
