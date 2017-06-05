package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type defaulEntityHandler struct {
	entity Entity
}

var v1Path = "/api/v1/"
var kProducer KafkaProducer
var kConsumer KafkaConsumer

func listen() error {
	p, err := newKafkaSyncProducer()
	if err != nil {
		return err
	}
	kProducer = KafkaProducer{make(chan *Event, 10), p}
	kProducer.start()

	http.Handle(v1Path+"costumer/", &defaulEntityHandler{&Costumer{}})
	http.Handle(v1Path+"account/", &defaulEntityHandler{&Account{}})
	http.Handle(v1Path+"transaction/", &defaulEntityHandler{&Transaction{}})
	if err := http.ListenAndServe(getPort(), nil); err != nil {
		return err
	}

	return nil
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
		if event, err := entity.create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			kProducer.c <- event
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(entity); err != nil {
			log.Println("SEVERE: %v error returning json response %v\n", err, entity)
		}
		return
	case "GET":
		a_path := strings.Split(r.URL.Path, "/")

		if "" != a_path[4] { //by id
			entity := h.entity.newEntity()
			if err := findOne(entity.collection(), entity, a_path[4]); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(entity); err != nil {
				log.Println("SEVERE: %v error returning json response %v\n", err, entity)
			}
		} else {

			entity := h.entity.newEntity()
			entities := reflect.New(reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(entity)), 0, 0).Type())

			if err := findAll(entity.collection(), entities.Interface()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(entities.Interface()); err != nil {
				log.Println("SEVERE: %v error returning json response %v\n", err, entity)
			}
		}
		return

	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`<!DOCTYPE html>
	<html>
	<head>
	<title>Page Title</title>
	</head>
	<body>
	
	<h1>litebank</h1>
	<p>Not so fast.</p>
	
	</body>
	</html>`))
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
