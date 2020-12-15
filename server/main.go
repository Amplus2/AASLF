package main

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Player struct {
	Name    string
	Session string
}

type Game struct {
	Players []Player
	Name    string
	ID      string
	Running bool
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

//TODO: learn how to do random in golang
func tempRand(b []byte) {
	b[8] = 4
	b[9] = 5
	b[10] = 6
	b[11] = 7
	b[12] = 8
	b[13] = 9
	b[14] = 10
	b[15] = 11
}

func GeneratePlayerSession() string {
	nonceMutex.Lock()
	var arr [16]byte
	binary.BigEndian.PutUint64(arr[:], nonce)
	nonce++
	nonceMutex.Unlock()
	tempRand(arr[:])
	var sum = sha1.Sum(arr[:])
	return base32.StdEncoding.EncodeToString(sum[:])
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var games []Game

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Use. POST. Requests.\"}")
		}

		var req struct {
			Game   string
			Player string
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Invalid JSON: "+err.Error()+"\"}")
		}

		game := Game{Name: req.Game, Players: []Player{Name: req.Player, Session: GeneratePlayerSession()}, ID: GenerateGameID()}

		games = append(games, game)

		var res struct {
			Status  string
			ID      string
			Session string
		}
		res.Status = "ok"
		res.ID = game.ID
		res.Session = game.Players[0].Session

		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Cannot encode JSON: "+err.Error()+"\"}")
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
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Invalid JSON: "+err.Error()+"\"}")
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
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Game not found.\"}")
		}

		session := GeneratePlayerSession()
		game.Players = append(game.Players, Player{Name: req.Player, Session: session})

		var res struct {
			Status  string
			Session string
		}
		res.Status = "ok"
		res.Session = session

		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Cannot encode JSON: "+err.Error()+"\"}")
		}
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})

	http.ListenAndServe(":1312", nil)
}
