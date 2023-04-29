package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/handlers"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/render"
	"github.com/alexedwards/scs/v2"
)


const portNumber = ":8000"
var app config.AppConfig
var sessions *scs.SessionManager

func main(){


	//what we are going to store in session
	gob.Register(models.Reservation{})

	//change this to true in production
	app.InProduction = false

	sessions = scs.New()
	sessions.Lifetime = 24 * time.Hour
	sessions.Cookie.Persist = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Secure = app.InProduction
	app.Session = sessions


	tc, err := render.CreateTemplateCache()
	if err != nil{
		log.Fatal("cannot create template cache.")
	}
	
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

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