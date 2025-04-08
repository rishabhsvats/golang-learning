package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)
	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate", app.ActivateAccount)
	mux.Get("/plan", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubscribeToPlan)
	mux.Get("/text-email", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "info@company.com",
			FromName:    "info",
			ErrorChan:   make(chan error),
		}
		msg := Message{
			To:      "me@test.com",
			Subject: "test email",
			Data:    "Hello World!",
		}
		m.sendMail(msg, make(chan error))
	})
	return mux
}
