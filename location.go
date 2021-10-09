package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type RegionID int

const (
	KantoRegionID RegionID = 1
)

type LocationArea struct {
	EncounterMethodRates []struct {
		EnconterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
}

type Location struct {
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

func (l Location) ExtractRegionIDFromURL() (int, error) {
	extract := strings.Builder{}
	for i := len(l.Region.URL) - 1; i > 0; i-- {
		char := l.Region.URL[i]
		if char == '/' {
			break
		}

		if char >= '0' && char <= '9' {
			extract.WriteByte(char)
		}
	}

	return strconv.Atoi(extract.String())
}

type PokemonLocation struct {
	EnconterMethods []string
	Location        string
	Area            struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location_area"`
}

func (pl *PokemonLocation) IsInRegion(regionID RegionID) (bool, error) {
	locationAreaDetails, err := http.Get(pl.Area.URL)
	if err != nil {
		fmt.Printf("failed to get location details err: %v\n", err)
		return false, err
	}

	defer locationAreaDetails.Body.Close()

	locationArea := LocationArea{}
	if err := json.NewDecoder(locationAreaDetails.Body).Decode(&locationArea); err != nil {
		fmt.Printf("failed to parse location details err: %v\n", err)
		return false, err
	}

	for _, encounterMethod := range locationArea.EncounterMethodRates {
		pl.EnconterMethods = append(
			pl.EnconterMethods,
			encounterMethod.EnconterMethod.Name,
		)
	}

	pl.Location = locationArea.Location.Name

	locationDetails, err := http.Get(locationArea.Location.URL)
	if err != nil {
		return false, err
	}

	defer locationDetails.Body.Close()

	location := Location{}
	if err := json.NewDecoder(locationDetails.Body).Decode(&location); err != nil {
		return false, err
	}

	if location.Region.Name == "kanto" {
		return true, nil
	}

	extractedRegionID, err := location.ExtractRegionIDFromURL()
	if err != nil {
		return false, err
	}

	if RegionID(extractedRegionID) == KantoRegionID {
		return true, nil
	}

	return false, nil
}
