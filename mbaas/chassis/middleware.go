package chassis

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// AddCommonMiddleware adds common middleware for all routes in all
// services. (This is intended to include *only* middleware that
// should be used both for internal services and those exposed
// externally. It's mostly about getting a consistent logging story to
// help with log aggregation.)
func AddCommonMiddleware(r chi.Router) {
	// Set up zerolog request logging.
	r.Use(hlog.NewHandler(log.Logger))
	logs := func(r *http.Request, status, size int, duration time.Duration) {
		basicRequestLog(r, status, size, duration).Msg("")
	}
	r.Use(hlog.AccessHandler(logs))
	r.Use(hlog.RemoteAddrHandler("ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))

	//r.Use(logHandler)

	// Panic recovery.
	r.Use(middleware.Recoverer)
}

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		x, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		save := r.Body
		if r.Body == nil {
			r.Body = nil
		} else {
			save, r.Body, _ = drainBody(r.Body)
		}
		reqBody, err := ioutil.ReadAll(save)
		if !utf8.Valid(reqBody) {
			reqBody = []byte("** BINARY DATA IN BODY **")
		}
		x = append(x, reqBody...)
		log.Info().Str("dir", "request").Msg(string(x))
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)
		resp := fmt.Sprintf("%d\n", rec.Code)
		for k, v := range rec.HeaderMap {
			resp += k + ": " + strings.Join(v, ",") + "\n"
		}
		body := rec.Body.String()
		if !utf8.Valid([]byte(body)) {
			body = "** BINARY DATA IN BODY **"
		}
		log.Info().Str("dir", "response").Msg(resp + body)

		// this copies the recorded response to the response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// Basic HTTP request logging.
func basicRequestLog(r *http.Request, status, size int, duration time.Duration) *zerolog.Event {
	if r.URL.Path == "/healthz" {
		return nil
	}
	return hlog.FromRequest(r).Info().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Int("status", status).
		Int("size", size).
		Dur("duration", duration)
}
