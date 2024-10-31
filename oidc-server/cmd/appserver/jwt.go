package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rishabhsvats/oidc-server/pkg/oidc"
)

// gets token from tokenUrl validating token with jwksUrl and returning token & claims
func getTokenFromCode(tokenUrl, jwksUrl, redirectUri, clientID, clientSecret, code string) (*jwt.Token, *jwt.RegisteredClaims, error) {
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

	var token oidc.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, nil, fmt.Errorf("unmasrhal token error: %s", err)

	}

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token.IDToken, claims, func(*jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("token parsing failed: %s", err)
	}

	return parsedToken, claims, nil
}
