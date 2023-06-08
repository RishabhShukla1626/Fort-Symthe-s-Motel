package main

import (
	"fmt"
	"testing"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T){
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type){
	case *chi.Mux: 
	//do nothing, test passed

	default: 
		t.Error(fmt.Sprintf("type is not a *chi.Mux, instead it's %T", v))
	}
}