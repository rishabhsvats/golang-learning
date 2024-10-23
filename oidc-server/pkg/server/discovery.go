package server

import (
	"encoding/json"
	"net/http"

	"github.com/rishabhsvats/oidc-server/pkg/oidc"
)

func (s *server) discovery(w http.ResponseWriter, r *http.Request) {
	discovery := oidc.Discovery{
		Issuer:                            s.Config.Url,
		AuthorizationEndpoint:             s.Config.Url + "/authorization",
		TokenEndpoint:                     s.Config.Url + "/token",
		UserinfoEndpoint:                  s.Config.Url + "/userinfo",
		JwksURI:                           s.Config.Url + "/jwks.json",
		ScopesSupported:                   []string{"oidc"},
		ResponseTypesSupported:            []string{"code"},
		TokenEndpointAuthMethodsSupported: []string{"none"},
	}
	out, err := json.Marshal(discovery)
	if err != nil {
		returnError(w, err)

	}
	w.Write(out)
}
