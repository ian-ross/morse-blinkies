package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ian-ross/morse-blinkies/mbaas/model"
	"github.com/rs/zerolog/log"
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

	text := r.FormValue("blinky-text")
	var rules *model.Rules
	ruleText := r.FormValue("blinky-rules")
	if ruleText == "" {
		var err error
		rules, err = rulesFromForm(r, s.baseRules)
		if err != nil {
			s.tmpl.ExecuteTemplate(w, "error", err.Error())
			return
		}
	} else {
		inRules := model.Rules{}
		if err := json.Unmarshal([]byte(ruleText), &inRules); err != nil {
			print(err)
			s.tmpl.ExecuteTemplate(w, "error", "Couldn't parse rules text!")
		}
		rules = &inRules
	}

	_, _, fullProjName, err := rules.ProjectName(text)
	if err != nil {
		s.tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}
	html := filepath.Join(s.bm.OutputDir, fullProjName+".html")
	if _, err := os.Stat(html); err == nil {
		http.Redirect(w, r, "/output/"+fullProjName+".html", http.StatusFound)
		return
	}

	s.queuer.Submit(fullProjName, text, rules)
	s.tmpl.ExecuteTemplate(w, "pending", fullProjName)
}

func (s *Server) jobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	conn, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't upgrade connection to WebSocket"))
		log.Error().Msg("couldn't upgrade connection to WebSocket")
		return
	}

	// Set up status channel.
	ch, exists := s.queuer.Subscribe(jobID)
	if !exists {
		fmt.Println("Oops: exists is false from Subscribe")
		// Job should already have been processed.
		conn.WriteJSON(URLNotification{"/output/" + jobID + ".html"})
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		close(ch)
		return nil
	})

	// Process messages until channel closed (triggered by WebSocket
	// closure).
	for msg := range ch {
		if err = conn.WriteJSON(msg); err != nil {
			log.Error().Msg("couldn't write status message to browser")
			return
		}
	}
}

func rulesFromForm(r *http.Request, baseRules *model.Rules) (*model.Rules, error) {
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

	case "single":
		// OK

	default:
		return nil, errors.New("invalid blinky type field")
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
		LEDForwardVoltage: baseRules.LEDForwardVoltage,
		LEDForwardCurrent: baseRules.LEDForwardCurrent,
		LEDGroups:         ledGroups,
		BlinkRate:         blinkRate,
		TransistorDrivers: transistors,
		Footprints:        baseRules.Footprints,
	}, nil
}
