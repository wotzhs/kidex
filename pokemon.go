package main

import (
	"encoding/json"
	"net/http"
)

type Type struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Types          []Type `json:"types"`
	Stats          []Stat `json:"stats"`
	PokeAPIBaseURL string `json:"-"`
}

func (p *Pokemon) Find(identifier string) error {
	pokemonDetails, err := http.Get(p.PokeAPIBaseURL + "/pokemon/" + identifier)
	if err != nil {
		return err
	}

	defer pokemonDetails.Body.Close()

	return json.NewDecoder(pokemonDetails.Body).Decode(&p)
}

func (p Pokemon) PrettyPrint() (string, error) {
	b, err := json.MarshalIndent(p, "", "\t")
	return string(b), err
}
