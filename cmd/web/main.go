package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/joho/godotenv"

	"github.com/alexedwards/scs/v2"
	"github.com/kiniconnect/booking-app/internals/config"
	"github.com/kiniconnect/booking-app/internals/drivers"
	"github.com/kiniconnect/booking-app/internals/handlers"
	"github.com/kiniconnect/booking-app/internals/helpers"
	"github.com/kiniconnect/booking-app/internals/models"
	"github.com/kiniconnect/booking-app/internals/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*drivers.DB, error) {

	var DATABASE_URL string
	if app.InProduction {
		DATABASE_URL = os.Getenv("DATABASE_URL")
		if DATABASE_URL == "" {
			log.Fatal("DATABASE_URL is not set in production environment")
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		DATABASE_URL = os.Getenv("DATABASE_URL")
		if DATABASE_URL == "" {
			log.Fatal("DATABASE_URL is not set")
		}
	}

	// what i want to put in session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})	
	gob.Register(models.Restriction{})	
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})	

	// change this to true when in production
	app.InProduction = false


	// create a custom logger
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to the database
	log.Println("Connecting to database...")
	db, err := drivers.ConnectToSQL(DATABASE_URL)
	if err != nil {
		log.Fatal("cannot connect to database")
	}

	fmt.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	log.Println("Application started successfully")

	return db, nil
}
