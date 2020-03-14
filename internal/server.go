package internal

import (
	"encoding/json"
	"mime"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	srv    *http.Server
	router *mux.Router
	db     *database
	Logger
}

func InitServer() *server {
	r := mux.NewRouter()
	r.Use(recoverMw, loggerMw)
	srv := &http.Server{
		Addr:        ":8080",
		Handler:     r,
		ReadTimeout: 1 << 20,
	}

	var db database
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     logrus.LevelHooks{},
		Level:     logrus.DebugLevel,
	}
	return &server{srv, r, &db, logger}
}

func (s *server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *server) Routes() {
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("im ok")) }).Methods("GET")
	s.router.HandleFunc("/api", s.handleVideo()).Methods("POST")
}

func (s *server) handleVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if SupportedSuffix(mediaType) {

		}
		var video video
		res, err := video.RawMeta(r.Body)
		if err != nil {
			s.Logger.Errorf("%v", err)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			s.Logger.Warnf("%v", err)
		}
	}
}

/**
 * Helpers
 */

type Logger interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func recoverMw(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func loggerMw(next http.Handler) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handlerFunc)
}
