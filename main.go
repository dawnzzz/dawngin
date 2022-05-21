package main

import (
	"DawnGin/dain"
	"fmt"
	"net/http"
)

func main() {
	e := dain.New()

	e.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World, URL path = %v", r.URL.Path)
	})

	e.Post("/", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	e.Run(":9000")
}
