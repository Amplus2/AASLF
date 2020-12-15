package main

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Game struct {
	Players []string
	Name    string
	ID      string
}

//not a nonce like in aes ctr,
//but a counter used for everything pseudo-random
var nonce uint64 = 0
var nonceMutex = &sync.Mutex{}

//maybe this should just be rand, but its fun this way
func GenerateGameID() string {
	nonceMutex.Lock()
	var arr [8]byte
	binary.BigEndian.PutUint64(arr[:], nonce)
	nonce++
	nonceMutex.Unlock()
	sum := sha1.Sum(arr[:])
	return strings.ToUpper(base32.StdEncoding.EncodeToString(sum[:5]))
}

//TODO: generate session

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

		game := Game{Name: req.Game, Players: []string{req.Player}, ID: GenerateGameID()}

		games = append(games, game)

		var res struct {
			Status string
			ID     string
			//TODO: session
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

		var game Game
		var found bool
		for _, g := range games {
			if g.ID == req.Game {
				game = g
				found = true
				break
			}
		}

		if !found {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Can't find the game")
		}

		game.Players = append(game.Players, req.Player)

		var res struct {
			Status string
			//TODO: session
		}
		res.Status = "ok"

		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Can't even encode JSON.")
		}
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})

	http.ListenAndServe(":1312", nil)
}
