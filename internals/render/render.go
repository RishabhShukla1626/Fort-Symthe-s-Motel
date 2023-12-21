package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/justinas/nosurf"
)

var tc = make(map[string]*template.Template)
var app *config.AppConfig

func NewTemplates(a *config.AppConfig){
	app = a
}


func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id"){
		td.IsAuthenticated = 1
	}
	return td
}


func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	//create template cache
	if app.UseCache{
		tc = app.TemplateCache
	} else{
		tc, _ = CreateTemplateCache()
	}
	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok{
		log.Fatal("could not get the template cache")
		return errors.New("can't get template from cache")
	}
	//render the template
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)
	
	_, err := buf.WriteTo(w)

	if err != nil{
		log.Println(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil{
		log.Println(err)
		return err
	}

	return nil

}


func CreateTemplateCache() (map[string]*template.Template, error) {
	mycache := map[string]*template.Template{}

	//get all the files with page.html
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil{
		return mycache, err
	}
	// range through all the files with *page.html
	for _, page := range pages{
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page) 
		if err != nil{
			return mycache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil{
			return mycache, err
		}
		if len(matches) > 0{
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil{
				return mycache, err
			}
		}
		mycache[name] = ts
	}
	return mycache, nil
}


//Simple Template cache
// func RenderTemplate(w http.ResponseWriter, t string){
// 	var tmpl *template.Template
// 	var err error

// 	_, inMap := tc[t]
// 	if !inMap{
// 		// need to create the template
// 		log.Println("creating template and adding to cache.")

// 		err = createTemplatecache(t)
		
// 		if err != nil{
// 			log.Println(err)
// 			}
// 	} else{
// 		//have template in cache
// 		log.Println("using cached template")
// 	}
	
// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil{
// 			log.Println(err)
// 		}
// }

// func createTemplatecache(t string) error {
// 	templates := []string{
// 		fmt.Sprint("./templates/"+t),
// 		"./templates/base.layout.html",
// 	}

// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil{
// 		return err
// 	}
// 	tc[t] = tmpl
// 	return nil
// }