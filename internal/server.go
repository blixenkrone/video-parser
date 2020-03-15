package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/blixenkrone/video-parser/encoder"
	"github.com/davecgh/go-spew/spew"
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
	s.router.HandleFunc("/test", s.testHandler()).Methods("POST")
}

func (s *server) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *server) handleVideo() http.HandlerFunc {
	type response struct {
		Meta      *encoder.FFMPEGMetaOutput `json:"meta,omitempty"`
		Thumbnail []byte                    `json:"thumbnail,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		supported := encoder.SupportedSuffix(mediaType)
		if !supported {
			http.Error(w, "not supported media", 500)
			return
		}
		var res response
		var buf bytes.Buffer
		tr := io.TeeReader(r.Body, &buf)
		defer r.Body.Close()

		ffmpegOut, err := encoder.RawMeta(tr)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting raw meta: %v", err), 500)
			return
		}
		res.Meta = ffmpegOut.SanitizeOutput()

		thumb, err := encoder.Thumbnail(&buf)
		if err != nil {
			return
		}
		spew.Dump("thumb")
		res.Thumbnail = thumb

		spew.Dump(thumb)

		if err := json.NewEncoder(w).Encode(&res); err != nil {
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
