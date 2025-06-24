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
	log.Println(input)
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

func (app *Application) AdminLogin(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Credential string `json:"credential"`
		Password   string `json:"password"`
	}
	var input Input
	err := app.readJSON(r, &input)
	if err != nil {
		log.Println("input error:", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Println(input)

	admin, err := app.AdminLoginFunc(input.Credential, input.Password)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	payload := struct {
		Admin   AdminUser `json:"admin_user"`
		Message string    `json:"message"`
	}{
		Admin:   admin,
		Message: "Login successful",
	}
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Application) Admin(w http.ResponseWriter, r *http.Request) {
	teams, projects, err := app.GetDataForAdmin()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	payload := struct {
		Teams    []Team            `json:"teams"`
		Projects []ProjectTeamJoin `json:"projects"`
	}{
		Teams:    teams,
		Projects: projects,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Application) ProjectSubmit(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		TeamName    string `json:"team_name"`
		ProjectRepo string `json:"project_repo"`
		ImageLink   string `json:"image_link"`
	}
	var input Input

	err := app.readJSON(r, &input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	project, message, err := app.SubmitProject(input.TeamName, input.ProjectRepo, input.ImageLink)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	payload := struct {
		Message string  `json:"message"`
		Project Project `json:"project"`
	}{
		Message: message,
		Project: project,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Application) Project(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		TeamName string `json:"team_name"`
	}
	var input Input

	err := app.readJSON(r, &input)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	project, err := app.GetProjects(input.TeamName)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := struct {
		Project Project `json:"project"`
	}{
		Project: project,
	}
	app.writeJSON(w, http.StatusOK, payload)
}
