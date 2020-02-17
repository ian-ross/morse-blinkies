package server

import (
	"net/http"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "home", nil)
}

func (s *Server) advanced(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "advanced", nil)
}

func (s *Server) newJob(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.tmpl.ExecuteTemplate(w, "error", "Couldn't parse form response!")
		return
	}

	// ruleText := r.FormValue("rules")
	// if ruleText == "" {
	text := r.FormValue("blinky-text")
	// ledType := r.FormValue("led-type")
	// blinkRate, err := strconv.Atoi(r.FormValue("blink-rate"))
	// if err != nil {
	// 	s.tmpl.ExecuteTemplate(w, "error", "Invalid blink rate value!")
	// 	return
	// }

	// fmt.Printf("TEXT: '%s'\n", text)
	// fmt.Printf("TYPE: %s\n", ledType)
	// fmt.Printf("RATE: %d\n", blinkRate)
	// }

	htmlURL, err := s.bm.Make(text, s.tmpl)
	if err != nil {
		s.tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}

	http.Redirect(w, r, htmlURL, http.StatusFound)
	// jobID := 1
	// s.tmpl.ExecuteTemplate(w, "pending", jobID)
}

func (s *Server) pendingJob(w http.ResponseWriter, r *http.Request) {

}
