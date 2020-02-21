package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ian-ross/morse-blinkies/mbaas/model"
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
	fmt.Println(r.Form)

	text := r.FormValue("blinky-text")
	var rules *model.Rules
	ruleText := r.FormValue("rules")
	if ruleText == "" {
		var err error
		rules, err = rulesFromForm(r)
		if err != nil {
			s.tmpl.ExecuteTemplate(w, "error", err.Error())
			return
		}
	} else {
		if err := json.Unmarshal([]byte(ruleText), rules); err != nil {
			s.tmpl.ExecuteTemplate(w, "error", "Couldn't parse rules text!")
		}
	}

	htmlURLOrError, err := s.bm.Make(text, rules, s.tmpl)
	if err != nil {
		msg := err.Error()
		if htmlURLOrError != "" {
			msg = htmlURLOrError
		}
		s.tmpl.ExecuteTemplate(w, "error", msg)
		return
	}

	http.Redirect(w, r, htmlURLOrError, http.StatusFound)
	// jobID := 1
	// s.tmpl.ExecuteTemplate(w, "pending", jobID)
}

func (s *Server) pendingJob(w http.ResponseWriter, r *http.Request) {

}

func rulesFromForm(r *http.Request) (*model.Rules, error) {
	ledType := r.FormValue("led-type")
	ledGroups := []int{}
	switch ledType {
	case "multi":
		led := 1
		for {
			cnts := r.FormValue(fmt.Sprintf("led-count-%d", led))
			if cnts == "" {
				break
			}
			cnt, err := strconv.Atoi(cnts)
			if err != nil {
				return nil, errors.New("invalid LED group size")
			}
			ledGroups = append(ledGroups, cnt)
			led++
		}

	case "group":
		cnts := r.FormValue("led-count")
		if cnts == "" {
			return nil, errors.New("missing LED group size")
		}
		cnt, err := strconv.Atoi(cnts)
		if err != nil {
			return nil, errors.New("invalid LED group size")
		}
		ledGroups = append(ledGroups, cnt)
	}

	blinkRate, err := strconv.Atoi(r.FormValue("blink-rate"))
	if err != nil {
		return nil, errors.New("invalid blink rate value")
	}

	transistors := false
	if r.FormValue("transistor-drivers") == "on" {
		transistors = true
	}

	return &model.Rules{
		Type:              model.BlinkyType(ledType),
		LEDForwardVoltage: 2.2,
		LEDForwardCurrent: 20.0,
		LEDGroups:         ledGroups,
		BlinkRate:         blinkRate,
		TransistorDrivers: transistors,
	}, nil
}
