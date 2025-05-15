package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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


	// what i want to put in session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})

	// Read flag
	inProduction := flag.Bool("inproduction", true, "Run the application in production mode")
	useCache := flag.Bool("cache", true, "Use cache for templates")
	 dbURL := flag.String("dburl", "", "Database connection URL")

	// Parse command line flags
	flag.Parse()

	if *dbURL == "" {
		fmt.Println("Database URL is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Using connection string:", *dbURL)


 
	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache  
	app.InProduction = *inProduction
	app.UseCache = *useCache

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
	 db, err := drivers.ConnectToSQL(*dbURL)
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
	 

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	log.Println("Application started successfully")

	return db, nil
}
