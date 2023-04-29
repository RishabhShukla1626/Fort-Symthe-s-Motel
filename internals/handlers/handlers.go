package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/forms"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/render"
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
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request){
	    remoteIp := m.App.Session.GetString(r.Context(), "remoteIp")
		fmt.Println("remote_ip:",remoteIp)
		render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request){
	
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}


//Reservation renders the make-reservation page and displays a form 
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request){
	
	var emptyReservation models.Reservation
	
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form : forms.New(nil),
		Data: data,
	})
}


//PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()

	if err != nil{
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid(){
			data := make(map[string]interface{})
			data["reservation"] = reservation
			
			render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form : form,
			Data: data,
		})
	return 
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}



func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request){
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	
	if !ok{
		log.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}


func (m *Repository) Availability(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}


type jsonResponse struct{
	OK bool `json: "ok"`
	Message string `json: "message"`
}

//AvailabilityJSON: handles request to check availability and returns JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request){
	response := jsonResponse{
		OK : true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(response, "", "     ")
	if err != nil{
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
}



func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request){
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", start, end)))
}


func (m *Repository) Contact(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}