package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rishabhsvats/oidc-server/pkg/oidc"
)

const redirectUri = "http://localhost:8081/callback"

type app struct {
}

func main() {

	a := app{}

	http.HandleFunc("/", a.index)
	http.HandleFunc("/callback", a.callback)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %s\n", err)
	}
}

func (a *app) index(w http.ResponseWriter, r *http.Request) {
	oidcEndpoint := os.Getenv("OIDC_ENDPOINT")
	discovery, err := oidc.ParseDiscovery(oidcEndpoint + "/.well-known/openid-configuration")
	if err != nil {
		returnError(w, fmt.Errorf("parse discovery error: %s", err))
		return
	}
	state, err := oidc.GetRandomString(64)
	if err != nil {
		returnError(w, fmt.Errorf("GetRandomString error: %s", err))
		return
	}
	authorizationURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%sscope=openid&response_type=code&state=%s", discovery.AuthorizationEndpoint, os.Getenv("CLIENT_ID"), redirectUri, state)
	w.Write([]byte(`
	<html>
	<body>
	<a href="` + authorizationURL + `"><button>Login</button></a>
	</body>
	</html>`))
}

func (a *app) callback(w http.ResponseWriter, r *http.Request) {

}
func returnError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	fmt.Printf("Error: %s\n", err)
}
