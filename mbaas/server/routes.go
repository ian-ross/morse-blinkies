package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/ian-ross/morse-blinkies/mbaas/chassis"
)

func (s *Server) routes(devMode bool, csrfSecret string, outputDir string) chi.Router {
	r := chi.NewRouter()

	chassis.AddCommonMiddleware(r)

	r.Get("/", s.home)
	r.Get("/advanced", s.advanced)
	r.Post("/", s.newJob)
	r.Get("/status/{id}", s.jobStatus)
	r.Get("/sparkline", s.sparkline)

	// Static files.
	FileServer(r, "/assets", "assets")
	FileServer(r, "/info", "info")
	FileServer(r, "/output", outputDir)

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a directory.
func FileServer(r chi.Router, path string, root string) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	absroot, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}
	fs := http.StripPrefix(path, http.FileServer(http.Dir(absroot)))

	if path != "/" && path[len(path)-1] != '/' {
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
