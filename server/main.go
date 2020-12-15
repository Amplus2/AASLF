package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Game struct {
	Players []string
	Name    string
	ID      string
}

func main() {
	var games []Game

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "This is an API based on POST-requests.")
		}

		var req struct {
			Game   string
			Player string
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//TODO: json
			fmt.Fprintf(w, "Can't parse your JSON.")
		}

		game := Game{Name: req.Game, Players: []string{req.Player}, ID: "TODO"}

		games = append(games, game)

		var res struct {
			Status string
			ID     string
		}
		res.Status = "ok"
		res.ID = "TODO"

		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Can't even encode JSON.")
		}
	})

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "This is an API based on POST-requests.")
		}

		var req struct {
			Game   string
			Player string
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//TODO: json
			fmt.Fprintf(w, "Can't parse your JSON.")
		}

		//TODO: do it
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})

	http.ListenAndServe(":1312", nil)
}
