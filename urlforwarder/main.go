package main

import (
	"flag"
	"fmt"
	"github.com/oyekanmiayo/gophercises/urlforwarder/urlforwarder"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	fileNamePtr := flag.String(
		"f",
		"",
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
	default:
		db := initBuntDB()
		defer db.Close()
		handler, _ = urlforwarder.DBHandler(db, http.DefaultServeMux)
	}

	fmt.Println("Starting the server on :8080")
	_ = http.ListenAndServe(":8080", handler)
}

func initBuntDB() *buntdb.DB {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(
		func(tx *buntdb.Tx) error {
			_, _, err := tx.Set("/yuwa", "https://tere-sagay.com/", nil)
			_, _, err = tx.Set("/ayo", "https://twitter.com/_alternatewolf", nil)
			return err
		},
	)

	if err != nil {
		log.Printf("Couldn't insert any value into db: %v", err)
	}

	return db
}
