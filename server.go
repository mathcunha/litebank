package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type defaulEntityHandler struct {
	entity Entity
}

var v1Path = "/api/v1/"

func listen() error {
	http.Handle(v1Path+"costumer", &defaulEntityHandler{&Costumer{}})
	http.Handle(v1Path+"account", &defaulEntityHandler{&Account{}})
	return http.ListenAndServe(getPort(), nil)
}

func (h *defaulEntityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		entity := h.entity.newEntity()
		if err := decoder.Decode(entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := entity.create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(entity); err != nil {
			log.Println("SEVERE: %v error returning json response %v\n", err, entity)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(nil)
	return
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
