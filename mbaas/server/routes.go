package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/ian-ross/morse-blinkies/mbaas/chassis"
)

func (s *Server) routes(devMode bool, csrfSecret string) chi.Router {
	r := chi.NewRouter()

	chassis.AddCommonMiddleware(r)

	r.Get("/", s.home)
	r.Get("/advanced", s.advanced)
	r.Post("/", s.newJob)
	r.Get("/pending/{id:[0-9]+}", s.pendingJob)

	// Static files.
	FileServer(r, "/assets", http.Dir("assets"))
	FileServer(r, "/info", http.Dir("info"))
	FileServer(r, "/jobs", http.Dir("jobs"))

	return r
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "home", nil)
}

func (s *Server) advanced(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "advanced", nil)
}

func (s *Server) newJob(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) pendingJob(w http.ResponseWriter, r *http.Request) {

}
