package main

import (
	"log"
	"net/http"

	"github.com/a-bouts/tiles-server/wind"
)

func main() {

	w := wind.InitWinds()

	router := InitServer(w)
	log.Println("Start server on port 8090")
	log.Fatal(http.ListenAndServe(":8090", router))
}
