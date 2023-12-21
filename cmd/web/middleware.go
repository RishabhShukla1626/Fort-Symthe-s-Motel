package main

import (
	"fmt"
	"net/http"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/handlers"
	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		fmt.Println("Hit the Page")
		next.ServeHTTP(w, r)
	})
}


//Nosurf adds CSRFToken protection to all POST requests
func NoSurf(next http.Handler) http.Handler{
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}


//SessionLoad saves and loads the sessions on every request
func SessionLoad(next http.Handler) http.Handler{
	return sessions.LoadAndSave(next) 
}


func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !handlers.Repo.IsAuthenticated(r){
			app.Session.Put(r.Context(), "error", "Login first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
