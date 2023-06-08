package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/driver"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/forms"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/render"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/repository"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)


type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}


var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) (*Repository){
	return &Repository{
		App : a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok{
		log.Println("error")
		return
	}
	
	room, err := m.DB.GetRoomByID(res.RoomID)
	log.Println(room)
	if err != nil {
		log.Fatal(err)
		return
	}
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	stringMap :=  make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	data := make(map[string]interface{})
	data["reservation"] = res
	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form : forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}


//PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request){
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok{
		log.Println("error")
		return
	}
	err := r.ParseForm()

	if err != nil{
		log.Println(err)
		return
	}


	// sd := r.Form.Get("start_date")
	// ed := r.Form.Get("end_date")

	// layout := "2006-01-02"

	// startDate, err := time.Parse(layout, sd)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// endDate, err := time.Parse(layout, ed)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	reservation.FirstName =  r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	// reservation := models.Reservations{
	// 	FirstName: r.Form.Get("first_name"),
	// 	LastName: r.Form.Get("last_name"),
	// 	Email: r.Form.Get("email"),
	// 	Phone: r.Form.Get("phone"),
	// 	StartDate: startDate,
	// 	EndDate: endDate,
	// 	RoomID: roomID,
	// }

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

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		log.Fatal(err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestrictions{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	
	if err != nil {
		log.Fatal(err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}



func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request){
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservations)
	
	if !ok{
		log.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}


func (m *Repository) Availability(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}


type jsonResponse struct{
	OK bool `json: "ok"`
	Message string `json: "message"`
	RoomID string `json: "room_id"`
    StartDate string `json: "start_date"`
	EndDate string `json: "end_date"`
}

//AvailabilityJSON: handles request to check availability and returns JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request){
	
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-02-01"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	response := jsonResponse{
		OK : available,
		Message: "",
		RoomID: strconv.Itoa(roomID),
		StartDate: sd,
		EndDate: ed,
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
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		log.Fatal(err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		log.Fatal(err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil{
		log.Fatal(err)
		return
	}

	if len(rooms) == 0{
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms

	res := models.Reservations{
		StartDate: startDate,
		EndDate: endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.RenderTemplate(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}


func (m *Repository) Contact(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}


//ChooseRoom: Displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request){
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservations)
	if !ok{
		log.Println(err)
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}



//BookRoom: takes URL param, builds a sessional  variable, and 
//takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request){
	// id, s, e
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-02-01"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservations
	room, err := m.DB.GetRoomByID(roomID)
	log.Println(room)
	if err != nil {
		log.Fatal(err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}