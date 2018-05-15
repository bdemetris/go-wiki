package main

import (
	"log"
	"net/http"

	"github.com/bdemetris/wiki/wiki"
)

func main() {

	wiki.Open()
	defer wiki.Close()

	http.HandleFunc("/", wiki.HomeHandler)
	http.HandleFunc("/view/", wiki.MakeHandler(wiki.ViewHandler))
	http.HandleFunc("/edit/", wiki.MakeHandler(wiki.EditHandler))
	http.HandleFunc("/save/", wiki.MakeHandler(wiki.SaveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
