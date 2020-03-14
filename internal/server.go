package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	srv    *http.Server
	router *mux.Router
	db     *database
	logger
}

type logger interface {
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

func InitServer() *server {
	r := mux.NewRouter()
	r.Use(recoverMw, loggerMw)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	var db database
	logger := logrus.New()
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
		supported := SupportedSuffix(mediaType)
		s.Infof("is supported: %v", supported)
		if !supported {
			http.Error(w, err.Error(), 500)
			return
		}
		var video *video
		mrd, err := video.RawMeta(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		s.Infof("hello")
		// defer r.Body.Close()
		// b, err := ioutil.ReadAll(mrd)
		// if err != nil {
		// 	http.Error(w, err.Error(), 500)
		// 	return
		// }
		if err := json.NewEncoder(w).Encode(mrd); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

/**
 * Middleware
 */

func recoverMw(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("%v\n", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func loggerMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("<< %s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

// https://github.com/shimberger/gohls/blob/7c2a1cc3a0874acae3528dacca399eef3630aa5c/internal/cmd/root.go
