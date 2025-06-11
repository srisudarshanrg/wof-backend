package app

import (
	"fmt"
	"log"
	"net/http"
)

func (app Application) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Application running...")
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	var input Team
	err := app.readJSON(r, &input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	loggedIn, team, message, err := app.LoginTeam(input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusAccepted)
		return
	}

	output := struct {
		LoggedIn bool   `json:"logged_in"`
		Team     Team   `json:"team"`
		Message  string `json:"message"`
	}{
		LoggedIn: loggedIn,
		Team:     team,
		Message:  message,
	}

	err = app.writeJSON(w, http.StatusOK, output)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	var input Team
	err := app.readJSON(r, &input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Println(input)

	message, err := app.AddTeam(input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	output := struct {
		Message string `json:"message`
	}{
		Message: message,
	}

	app.writeJSON(w, http.StatusOK, output)
}
