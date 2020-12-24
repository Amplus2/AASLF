package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type player struct {
	Name  string
	Admin bool
}

type gameStatus int

const (
	lobby gameStatus = iota
	running
	voting
)

type game struct {
	Players    []player
	Name       string
	Categories []string
	ID         string
	Status     gameStatus
}

var games []game
var sessions map[string]map[string]string

func generateGameID() string {
	var arr [5]byte
	//TODO: err handling for this
	rand.Read(arr[:])
	return strings.ToUpper(base32.StdEncoding.EncodeToString(arr[:]))
}

func generatePlayerSession() string {
	var arr [20]byte
	//TODO: err handling for this
	rand.Read(arr[:])
	return base32.StdEncoding.EncodeToString(arr[:])
}

func httpBadRequest(w http.ResponseWriter, msg string) {
	httpError(w, http.StatusBadRequest, msg)
}

func httpError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"%s\"}\n", msg)
}

func beginPostHandler(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if r.Method != http.MethodPost {
		httpBadRequest(w, "Use. POST. Requests.")
		return true
	}

	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		httpBadRequest(w, "Invalid JSON: "+err.Error())
		return true
	}

	return false
}

func endHTTPHandler(w http.ResponseWriter, res interface{}) {
	err := json.NewEncoder(w).Encode(res)

	if err != nil {
		httpError(w, http.StatusInternalServerError, "Cannot encode JSON: "+err.Error())
	}
}

func searchGame(id string) (game, bool) {
	var game game
	var found bool
	for _, g := range games {
		if g.ID == id {
			game = g
			found = true
			break
		}
	}
	return game, found
}

func searchPlayer(game game, name string, session string) (player, bool, bool) {
	var player player
	found := false
	valid := false
	for _, p := range game.Players {
		if p.Name == name {
			player = p
			found = true
			valid = sessions[game.ID][p.Name] == session
			break
		}
	}
	return player, found, valid
}

func main() {
	http.HandleFunc("/v1/new", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game       string
			Player     string
			Categories []string
		}

		if beginPostHandler(w, r, &req) {
			return
		}

		game := game{Name: req.Game, Players: []player{{Name: req.Player, Admin: true}}, ID: generateGameID(), Categories: req.Categories, Status: lobby}
		games = append(games, game)
		sessions[game.ID][req.Player] = generatePlayerSession()

		var res struct {
			Status  string
			ID      string
			Session string
		}
		res.Status = "ok"
		res.ID = game.ID
		res.Session = sessions[game.ID][req.Player]

		endHTTPHandler(w, res)
	})

	http.HandleFunc("/v1/join", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game   string
			Player string
		}

		if beginPostHandler(w, r, &req) {
			return
		}

		game, found := searchGame(req.Game)

		if !found {
			httpBadRequest(w, "Game not found.")
			return
		}

		_, playerAlreadyExists, _ := searchPlayer(game, req.Player, "")

		if playerAlreadyExists {
			httpBadRequest(w, "That player name is already in use.")
			return
		}

		game.Players = append(game.Players, player{Name: req.Player})
		sessions[game.ID][req.Player] = generatePlayerSession()

		var res struct {
			Status  string
			Session string
		}
		res.Status = "ok"
		res.Session = sessions[game.ID][req.Player]

		endHTTPHandler(w, res)
	})

	http.HandleFunc("/v1/start", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game    string
			Player  string
			Session string
		}

		if beginPostHandler(w, r, &req) {
			return
		}

		game, found := searchGame(req.Game)

		if !found {
			httpBadRequest(w, "Game not found.")
			return
		}

		var player player
		var valid bool
		player, found, valid = searchPlayer(game, req.Player, req.Session)

		if !found {
			httpBadRequest(w, "Player not found.")
			return
		}

		if !valid {
			httpBadRequest(w, "Invalid session.")
			return
		}

		if !player.Admin {
			httpBadRequest(w, "No permission.")
			return
		}

		if game.Status != lobby {
			httpBadRequest(w, "That game is already running.")
			return
		}

		game.Status = running

		var res struct {
			Status string
		}
		res.Status = "ok"

		endHTTPHandler(w, res)
	})

	http.HandleFunc("/v1/stop", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game    string
			Player  string
			Session string
		}

		if beginPostHandler(w, r, &req) {
			return
		}

		game, found := searchGame(req.Game)

		if !found {
			httpBadRequest(w, "Game not found.")
			return
		}

		var valid bool
		_, found, valid = searchPlayer(game, req.Player, req.Session)

		if !found {
			httpBadRequest(w, "Player not found.")
			return
		}

		if !valid {
			httpBadRequest(w, "Invalid session.")
			return
		}

		if game.Status != running {
			httpBadRequest(w, "That game is not running.")
			return
		}

		game.Status = voting

		var res struct {
			Status string
		}
		res.Status = "ok"

		endHTTPHandler(w, res)
	})

	http.HandleFunc("/v1/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/v1/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpBadRequest(w, "This is a GET endpoint.")
			return
		}

		var res struct {
			Status string
			Games  []game
		}
		res.Status = "ok"
		res.Games = games

		endHTTPHandler(w, res)
	})

	http.ListenAndServe(":1312", nil)
}
