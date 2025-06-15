package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wings_of_fire/app"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbPassword := os.Getenv("DATABASE_PASSWORD")

	log.Println("Port:", port)

	var app app.Application

	app.Deployed = false
	app.DatabaseDSN = fmt.Sprintf("host=postgresql-raptor.alwaysdata.net dbname=raptor_wings_of_fire port=5432 user=raptor password=%s", dbPassword)
	app.FrontendLink = "https://spark-hack-website.vercel.app"

	if app.Deployed {
		app.Port = fmt.Sprintf("0.0.0.0:%s", port)
	} else {
		app.Port = fmt.Sprintf("localhost:%s", port)
	}

	log.Println("Connecting to database")
	db, err := app.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected to database")
	app.DB = db
	defer db.Close()

	log.Println("Application running")
	http.ListenAndServe(app.Port, app.Routes())
}
