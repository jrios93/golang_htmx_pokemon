package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//go:embed views/*
var views embed.FS
var pokemons []Pokemon

func main() {
	t := template.Must(template.ParseFS(views, "views/*.html"))

	router := http.NewServeMux()

	// Página principal
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Solo al inicio, cargamos a Pichu
		if len(pokemons) == 0 {
			resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/pichu")
			if err != nil {
				http.Error(w, "Unable to grab the pokemon data", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			var data Pokemon
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				http.Error(w, "Unable to parse the Pokemon data", http.StatusInternalServerError)
				return
			}
			pokemons = append(pokemons, data)
		}

		if err := t.ExecuteTemplate(w, "index.html", pokemons); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	// Ruta para agregar nuevos pokémon
	router.HandleFunc("POST /poke", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
			return
		}

		name := strings.ToLower(r.FormValue("pokemon"))
		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
		if err != nil {
			http.Error(w, "Unable to fetch new Pokemon", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var data Pokemon
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse Pokemon data", http.StatusInternalServerError)
			return
		}

		pokemons = append(pokemons, data)

		// Renderizamos solo la tarjeta
		if err := t.ExecuteTemplate(w, "card.html", data); err != nil {
			http.Error(w, "Something went wrong rendering the card", http.StatusInternalServerError)
		}
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", router)
}
