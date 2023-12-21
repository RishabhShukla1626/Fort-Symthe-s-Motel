package config

import (
	"html/template"
	"log"

	"github.com/TheDevCarnage/FortSmythesMotel/internals/models"
	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache bool
	TemplateCache map[string]*template.Template
	log *log.Logger
	InProduction bool
	Session *scs.SessionManager
	MailChan chan models.MailData
}