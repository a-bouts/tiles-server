package main

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-bouts/tiles-server/wind"
	"github.com/gorilla/mux"
)

type server struct {
	w *wind.Winds
}

func InitServer(w *wind.Winds) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	s := server{w: w}

	router.HandleFunc("/tiles/-/healthz", s.healthz).Methods(http.MethodGet)
	router.HandleFunc("/tiles/wind/{t}/{z}/{x}/{y}", s.getWindTile).Methods(http.MethodGet)

	return router
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	type health struct {
		Status string `json:"status"`
	}

	json.NewEncoder(w).Encode(health{Status: "Ok"})
}

func (s *server) getWindTile(w http.ResponseWriter, r *http.Request) {
	t := mux.Vars(r)["t"]
	z, err := strconv.Atoi(mux.Vars(r)["z"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	x, err := strconv.Atoi(mux.Vars(r)["y"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(mux.Vars(r)["x"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(t, z, x, y)

	m, err := time.Parse("200601021504", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tile := wind.GenerateTile(s.w, z, x, y, m)

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, tile, nil); err != nil {
		log.Println("unable to encode image.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
