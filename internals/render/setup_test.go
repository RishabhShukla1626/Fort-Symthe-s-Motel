package render

import (
	"encoding/gob"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/alexedwards/scs/v2"

	// "log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// testApp.InfoLog = infoLog

	// errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// testApp.ErrorLog = errorLog

	// change this to true when in production
	testApp.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}


type responseWriter struct{

}


func (tw responseWriter) Header() http.Header{
	var h http.Header
	return h
}

func (tw responseWriter) WriteHeader(i int){

}

func (tw responseWriter) Write(b []byte) (int, error){
	length := len(b)
	return length, nil
}