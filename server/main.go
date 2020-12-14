package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "container/list"
)

type Game struct {
    Players *list.List
    Name string
    Id string
}

func main() {
    var games = list.New()

    http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "This is an API based on POST-requests.")
        }

        var nr struct {
            Game string
            Player string
        }

        json.NewDecoder(r.Body).Decode(&nr)

        game := games.PushBack(Game{Name: nr.Game, Players: list.New()}).Value

        game.Players.PushBack(nr.Player)
    })

    http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {})

    http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {})

    http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {})

    http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {})

    http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {})

    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})
}
