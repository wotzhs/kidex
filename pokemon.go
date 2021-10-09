package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Type struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
	}
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Pokemon struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Types          []Type      `json:"types"`
	Stats          []Stat      `json:"stats"`
	Encounters     interface{} `json:"encounters"`
	PokeAPIBaseURL string      `json:"-"`
}

type PokeomonEncounterLocations struct {
	Methods  []string `json:"methods"`
	Location string   `json:"location"`
}

func (p *Pokemon) Find(identifier string) error {
	pokemonDetails, err := http.Get(p.PokeAPIBaseURL + "/pokemon/" + identifier)
	if err != nil {
		return err
	}

	if pokemonDetails.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Pokemon %v not found", identifier)
	}

	if pokemonDetails.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", http.StatusText(pokemonDetails.StatusCode))
	}

	defer pokemonDetails.Body.Close()

	if err := json.NewDecoder(pokemonDetails.Body).Decode(&p); err != nil {
		return err
	}

	if err := p.listEncounterLocations(); err != nil {
		return err
	}

	return nil
}

func (p *Pokemon) listEncounterLocations() error {
	encounters, err := http.Get(p.PokeAPIBaseURL + "/pokemon/" + strconv.Itoa(p.ID) + "/encounters")
	if err != nil {
		return fmt.Errorf("failed to list encounter areas err: %v", err)
	}

	if encounters.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", http.StatusText(encounters.StatusCode))
	}

	defer encounters.Body.Close()

	pokemonLocations := []PokemonLocation{}
	if err := json.NewDecoder(encounters.Body).Decode(&pokemonLocations); err != nil {
		return fmt.Errorf("failed to parse encounter areas err: %v", err)
	}

	wg := sync.WaitGroup{}

	locations := []PokeomonEncounterLocations{}

	for _, pokemonLocation := range pokemonLocations {
		wg.Add(1)
		go func(pl PokemonLocation) {
			yes, err := pl.IsInRegion(KantoRegionID)
			if err == nil && yes {
				locations = append(locations, PokeomonEncounterLocations{
					Methods:  pl.EnconterMethods,
					Location: pl.Location,
				})
			}
			wg.Done()
		}(pokemonLocation)
	}

	wg.Wait()

	if len(locations) < 1 {
		p.Encounters = "-"
	} else {
		p.Encounters = locations
	}

	return nil
}

func (p Pokemon) PrettyPrint() (string, error) {
	b, err := json.MarshalIndent(p, "", "\t")
	return string(b), err
}
