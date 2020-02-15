package server

import (
	"html/template"

	"github.com/ian-ross/morse-blinkies/mbaas/chassis"
)

// Server is the server structure for the user service.
type Server struct {
	chassis.Server
	tmpl *template.Template
}

// NewServer creates the server structure for the user service.
func NewServer(cfg *Config) *Server {
	// Common server initialisation.
	s := &Server{}
	s.Init(cfg.Port, s.routes(cfg.DevMode, cfg.CSRFSecret))
	s.tmpl = template.Must(template.ParseGlob("templates/*.html"))

	return s
}
