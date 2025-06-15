package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/admin", app.Admin)

	mux.Post("/login", app.Login)
	mux.Post("/register", app.Register)
	mux.Post("/admin-login", app.AdminLogin)

	return mux
}
