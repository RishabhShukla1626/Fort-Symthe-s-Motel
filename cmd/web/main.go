package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/driver"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/handlers"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/render"
	"github.com/alexedwards/scs/v2"
)


const portNumber = ":8000"
var app config.AppConfig
var sessions *scs.SessionManager

func main(){


	db, err := run()
	if err != nil{
		log.Fatal(err)
	}
	defer db.SQL.Close()
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// fmt.Println(fmt.Sprintf("Starting the Application at port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	fmt.Println(fmt.Sprintf("Starting the Application at port %s", portNumber))
	
	serve := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err!= nil{
		log.Fatal(err)
	}
	
}


func run() (*driver.DB, error) {
	
	//what we are going to store in session
	gob.Register(models.Reservations{})
	gob.Register(models.Users{})
	gob.Register(models.Restrictions{})
	gob.Register(models.Rooms{})
	gob.Register(models.RoomRestrictions{})


	//change this to true in production
	app.InProduction = false

	sessions = scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Cookie.Persist = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Secure = app.InProduction
	app.Session = sessions

	//connect to the database
	log.Println("Connecting to the database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")

	if err != nil {
		log.Fatal("Cannot connect to the database! Dying...")
	}


	tc, err := render.CreateTemplateCache()
	if err != nil{
		log.Fatal("cannot create template cache.")
		return nil, err
	}
	
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return db, nil
}