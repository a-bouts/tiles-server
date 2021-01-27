package main

import (
	"log"
	"net/http"

	"github.com/a-bouts/tiles-server/wind"
)

func main() {

	p := wind.InitWinds()

	router := InitServer(p)
	log.Println("Start server on port 8091")
	log.Fatal(http.ListenAndServe(":8091", router))
}
