package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Game struct {
	Players []string
	Name    string
	Id      string
}

func main() {
	var games []Game

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "This is an API based on POST-requests.")
		}

		var nr struct {
			Game   string
			Player string
		}

		json.NewDecoder(r.Body).Decode(&nr)

		game := Game{Name: nr.Game, Players: []string{nr.Player}}

		games = append(games, game)
	})

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "This is an API based on POST-requests.")
		}

		var jr struct {
			Game   string
			Player string
		}

		json.NewDecoder(r.Body).Decode(&jr)

		//TODO: do it
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})
}
