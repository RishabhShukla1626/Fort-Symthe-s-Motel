package handlers

import (
	"fmt"
	"net/http"

	"github.com/TheDevCarnage/FortSmythesMotel/pkg/config"
	"github.com/TheDevCarnage/FortSmythesMotel/pkg/models"
	"github.com/TheDevCarnage/FortSmythesMotel/pkg/render"
)


type Repository struct {
	App *config.AppConfig
}


var Repo *Repository

func NewRepo(a *config.AppConfig) (*Repository){
	return &Repository{
		App : a,
	}
}

func NewHandlers(r *Repository){
	Repo = r
}


func (m *Repository) Home(w http.ResponseWriter, r *http.Request){
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remoteIp", remoteIp)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request){
	    remoteIp := m.App.Session.GetString(r.Context(), "remoteIp")
		fmt.Println("remote_ip:",remoteIp)
		render.RenderTemplate(w, "about.page.html", &models.TemplateData{
			StringMap: map[string]string{"test":"hello From Backend data"},
		})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request){
	
	render.RenderTemplate(w, "majors.page.html", &models.TemplateData{})
}


func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{})
}