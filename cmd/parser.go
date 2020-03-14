package main

import (
	"github.com/blixenkrone/video-parser/internal"
)

func main() {

	srv := internal.InitServer()
	srv.Routes()

	err := srv.ListenAndServe()
	if err != nil {
		srv.Fatalf("%", err)
	}

}
