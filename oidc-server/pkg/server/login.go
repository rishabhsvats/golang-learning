package server

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rishabhsvats/oidc-server/pkg/oidc"
	"github.com/rishabhsvats/oidc-server/pkg/users"
)

//go:embed templates/*
var templateFs embed.FS

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			returnError(w, fmt.Errorf("parse form error: %s", err))
			return
		}
		sessionID := r.PostForm.Get("sessionID")
		loginRequest, ok := s.LoginRequest[sessionID]
		if !ok {
			returnError(w, fmt.Errorf("session not found"))
			return
		}
		auth, user, err := users.Auth(r.PostForm.Get("login"), r.PostForm.Get("password"), "")
		if err != nil {
			returnError(w, fmt.Errorf("authentication error: %s", err))
			return
		}
		if auth {
			code, err := oidc.GetRandomString(64)
			if err != nil {
				returnError(w, fmt.Errorf("get random string error: %s", err))
				return
			}
			loginRequest.CodeIssuedAt = time.Now()
			loginRequest.User = user
			s.Codes[code] = loginRequest
			delete(s.LoginRequest, sessionID)
			w.Header().Add("location", fmt.Sprintf("%s?code=%s&state=%s", loginRequest.RedirectURI, code, loginRequest.State))
			w.WriteHeader(http.StatusFound)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication failed"))
		}
		return
	}
	var (
		sessionID string
	)
	if sessionID = r.URL.Query().Get("sessionID"); sessionID == "" {
		returnError(w, fmt.Errorf("sessionID is empty"))
		return
	}
	templateFile, err := templateFs.Open("templates/login.html")
	if err != nil {
		returnError(w, fmt.Errorf("templateFS open error: %s", err))
		return
	}
	templateFileBytes, err := io.ReadAll(templateFile)
	if err != nil {
		returnError(w, fmt.Errorf("readall error: %s", err))
		return
	}
	templateFileStr := strings.Replace(string(templateFileBytes), "$SESSIONID", sessionID, -1)
	w.Write([]byte(templateFileStr))
}
