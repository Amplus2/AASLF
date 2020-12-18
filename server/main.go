package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Player struct {
	Name  string
	Admin bool
}

type Game struct {
	Players    []Player
	Name       string
	Categories []string
	ID         string
	Running    bool
}

var games []Game
var sessions map[string]map[string]string

func GenerateGameID() string {
	var arr [5]byte
	//TODO: err handling for this
	rand.Read(arr[:])
	return strings.ToUpper(base32.StdEncoding.EncodeToString(arr[:]))
}

func GeneratePlayerSession() string {
	var arr [20]byte
	//TODO: err handling for this
	rand.Read(arr[:])
	return base32.StdEncoding.EncodeToString(arr[:])
}

func BadHttpRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, msg)
}

func BeginPostHandler(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if r.Method != http.MethodPost {
		BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Use. POST. Requests.\"}")
		return true
	}

	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Invalid JSON: "+err.Error()+"\"}")
		return true
	}

	return false
}

func EndHttpHandler(w http.ResponseWriter, res interface{}) {
	err := json.NewEncoder(w).Encode(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"Status\":\"err\",\"Msg\":\"Cannot encode JSON: "+err.Error()+"\"}")
	}
}

func SearchGame(id string) (Game, bool) {
	var game Game
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

func SearchPlayer(game Game, name string, session string) (Player, bool, bool) {
	var player Player
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

		if BeginPostHandler(w, r, &req) {
			return
		}

		game := Game{Name: req.Game, Players: []Player{{Name: req.Player, Admin: true}}, ID: GenerateGameID(), Categories: req.Categories}
		games = append(games, game)
		sessions[game.ID][req.Player] = GeneratePlayerSession()

		var res struct {
			Status  string
			ID      string
			Session string
		}
		res.Status = "ok"
		res.ID = game.ID
		res.Session = sessions[game.ID][req.Player]

		EndHttpHandler(w, res)
	})

	http.HandleFunc("/v1/join", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game   string
			Player string
		}

		if BeginPostHandler(w, r, &req) {
			return
		}

		game, found := SearchGame(req.Game)

		if !found {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Game not found.\"}")
			return
		}

		_, playerAlreadyExists, _ := SearchPlayer(game, req.Player, "")

		if playerAlreadyExists {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"That player name is already used.\"}")
			return
		}

		game.Players = append(game.Players, Player{Name: req.Player})
		sessions[game.ID][req.Player] = GeneratePlayerSession()

		var res struct {
			Status  string
			Session string
		}
		res.Status = "ok"
		res.Session = sessions[game.ID][req.Player]

		EndHttpHandler(w, res)
	})

	http.HandleFunc("/v1/start", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Game    string
			Player  string
			Session string
		}

		if BeginPostHandler(w, r, &req) {
			return
		}

		game, found := SearchGame(req.Game)

		if !found {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Game not found.\"}")
			return
		}

		var player Player
		var valid bool
		player, found, valid = SearchPlayer(game, req.Player, req.Session)

		if !found {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Player not found.\"}")
			return
		}

		if !valid {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"Invalid session.\"}")
			return
		}

		if !player.Admin {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"No permission.\"}")
			return
		}

		if game.Running {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"That game is already running.\"}")
			return
		}

		game.Running = true

		var res struct {
			Status string
		}
		res.Status = "ok"

		EndHttpHandler(w, res)
	})

	http.HandleFunc("/v1/stop", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/v1/submit", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/v1/vote", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			BadHttpRequest(w, "{\"Status\":\"err\",\"Msg\":\"This is a GET endpoint.\"}")
			return
		}

		var res struct {
			Status string
			Games  []Game
		}
		res.Status = "ok"
		res.Games = games

		EndHttpHandler(w, res)
	})

	http.ListenAndServe(":1312", nil)
}
