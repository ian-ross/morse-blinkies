package chassis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

// Server represents the shared functionality used in all backend REST
// servers.
type Server struct {
	Ctx      context.Context
	Srv      *http.Server
	muAtExit sync.Mutex
	atExit   []func()
}

// Init initialises all the common infrastructure used by REST
// servers.
func (s *Server) Init(port int, r chi.Router) {
	// Randomise ID generation.
	rand.Seed(int64(time.Now().Nanosecond()))

	s.Ctx = context.Background()
	s.atExit = []func(){}
	s.Srv = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", port),
	}
}

// Serve runs a server event loop.
func (s *Server) Serve() {
	errChan := make(chan error, 0)
	go func() {
		log.Info().
			Str("address", s.Srv.Addr).
			Msg("server started")
		err := s.Srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	signalCh := make(chan os.Signal, 0)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case <-signalCh:
	case err = <-errChan:
	}

	s.shutdown()
	s.runAtShutdown()

	if err == nil {
		log.Info().Msg("server shutting down")
	} else {
		log.Fatal().Err(err).Msg("server failed")
	}
}

// AddAtExit adds an exit handler function.
func (s *Server) AddAtExit(fn func()) {
	s.muAtExit.Lock()
	defer s.muAtExit.Unlock()
	s.atExit = append(s.atExit, fn)
}

// Shut down server.
func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(s.Ctx, 10*time.Second)
	defer cancel()
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Error().Err(err)
	}
}

// Run at-exit processing.
func (s *Server) runAtShutdown() {
	s.muAtExit.Lock()
	defer s.muAtExit.Unlock()
	for _, fn := range s.atExit {
		fn()
	}
}

// SimpleHandlerFunc is a HTTP handler function that signals internal
// errors by returning a normal Go error, and when successful returns
// a response body to be marshalled to JSON. It can be wrapped in the
// SimpleHandler middleware to produce a normal HTTP handler function.
type SimpleHandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

// SimpleHandler wraps a simpleHandler-style HTTP handler function to
// turn it into a normal HTTP handler function. Go errors from the
// inner handler are returned to the caller as "500 Internal Server
// Error" responses. Returns from successful processing by the inner
// handler and marshalled into a JSON response body.
func SimpleHandler(inner SimpleHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Run internal handler: returns a marshalable result and an
		// error, either of which may be nil.
		result, err := inner(w, r)

		// Propagate Go errors as "500 Internal Server Error" responses.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
			return
		}

		// No response body, so internal handler dealt with response
		// setup.
		if result == nil {
			return
		}

		// Marshal JSON response body.
		body, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(body)
	}
}
