package server

import (
	"html/template"

	"github.com/ian-ross/morse-blinkies/mbaas/chassis"
	"github.com/ian-ross/morse-blinkies/mbaas/model"
	"github.com/ian-ross/morse-blinkies/mbaas/processing"
)

// Server is the server structure for the user service.
type Server struct {
	chassis.Server
	tmpl      *template.Template
	bm        *processing.BlinkyMaker
	baseRules *model.Rules
	queuer    *Queuer
}

// NewServer creates the server structure for the user service.
func NewServer(cfg *Config) *Server {
	// Common server initialisation.
	s := &Server{}
	s.Init(cfg.Port, s.routes(cfg.DevMode, cfg.CSRFSecret, cfg.OutputDir))
	s.tmpl = template.Must(template.ParseGlob("templates/*.html"))
	s.baseRules = model.ReadBaseRules(cfg.LibDir)
	s.bm = &processing.BlinkyMaker{
		cfg.BlinkyScript,
		cfg.WorkDir,
		cfg.OutputDir,
		cfg.TemplateDir,
	}
	s.queuer = MakeQueuer(s.bm, s.tmpl)

	return s
}
