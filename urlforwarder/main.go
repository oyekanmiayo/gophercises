package main

import (
	"flag"
	"fmt"
	"github.com/oyekanmiayo/gophercises/urlforwarder/urlforwarder"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func main() {
	fileNamePtr := flag.String(
		"f",
		"redirection.yaml",
		"File that contains paths",
	)

	flag.Parse()

	http.HandleFunc(
		"/",
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = fmt.Fprintln(writer, "Hello")
		},
	)

	content, _ := ioutil.ReadFile(*fileNamePtr)

	ext := filepath.Ext(*fileNamePtr)

	var handler http.Handler
	switch ext {
	case ".yml":
	case ".yaml":
		handler, _ = urlforwarder.YAMLHandler(content, http.DefaultServeMux)
	case ".json":
		handler, _ = urlforwarder.JSONHandler(content, http.DefaultServeMux)
	}

	fmt.Println("Starting the server on :8080")
	_ = http.ListenAndServe(":8080", handler)
}
