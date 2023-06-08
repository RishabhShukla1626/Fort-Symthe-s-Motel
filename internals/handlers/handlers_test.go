package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "smith"},
		{key: "email", value: "john@gmail.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	
	defer testServer.Close()

	for _, e := range theTests{
		if e.method=="GET"{
			//create client and call get method with testservaers port and append the main url
			resp, err := testServer.Client().Get(testServer.URL + e.url)
			
			if err != nil{
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s, expected %d but got %d instead", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else{
			values := url.Values{}
			for _, x := range e.params{
				values.Add(x.key, x.value)
			}

			resp, err := testServer.Client().PostForm(testServer.URL+e.url, values)

			if err != nil{
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s, expected %d but got %d instead", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}