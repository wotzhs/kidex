package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	CSVHeader []string = []string{"id", "name", "types", "stats", "encounters", "cached_at"}
)

type CacheEntry struct {
	Pokemon  Pokemon
	CachedAt time.Time
}

type Cache struct {
	NameMap map[string]string
	IDMap   map[string]CacheEntry
}

const (
	pos_ID int = iota
	pos_NAME
	pos_TYPES
	pos_STATS
	pos_ENCOUNTERS
	pos_CACHED_AT
)

func (c *Cache) Restore() error {
	// prevent wrongful restore
	if c.NameMap != nil && c.IDMap != nil {
		return nil
	}

	c.NameMap = map[string]string{}
	c.IDMap = map[string]CacheEntry{}

	f, err := os.OpenFile("cache", os.O_RDONLY, 0444)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)

	// skip header
	r.Read()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			continue
		}

		ID, err := strconv.Atoi(record[pos_ID])
		if err != nil {
			continue
		}

		types := []Type{}
		if err := json.Unmarshal([]byte(record[pos_TYPES]), &types); err != nil {
			continue
		}

		stats := []Stat{}
		if err := json.Unmarshal([]byte(record[pos_STATS]), &stats); err != nil {
			continue
		}

		pokemon := Pokemon{
			ID:         ID,
			Name:       record[pos_NAME],
			Types:      types,
			Stats:      stats,
			Encounters: "-",
		}

		encounters := []PokeomonEncounterLocations{}
		if err := json.Unmarshal([]byte(record[pos_ENCOUNTERS]), &encounters); err == nil {
			pokemon.Encounters = encounters
		}

		cachedAt, err := time.Parse(time.RFC3339, record[pos_CACHED_AT])
		if err != nil {
			continue
		}

		c.NameMap[pokemon.Name] = record[pos_ID]
		c.IDMap[record[pos_ID]] = CacheEntry{
			Pokemon:  pokemon,
			CachedAt: cachedAt,
		}
	}

	return nil
}

func (c Cache) WriteToCSV() error {
	f, err := os.OpenFile("cache", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)

	records := [][]string{CSVHeader}
	for name, id := range c.NameMap {
		cacheEntry, ok := c.IDMap[id]
		if !ok {
			continue
		}

		var (
			typesBuffer      bytes.Buffer
			statsBuffer      bytes.Buffer
			encountersBuffer bytes.Buffer
		)

		if err := json.NewEncoder(&typesBuffer).Encode(cacheEntry.Pokemon.Types); err != nil {
			continue
		}

		if err := json.NewEncoder(&statsBuffer).Encode(cacheEntry.Pokemon.Stats); err != nil {
			continue
		}

		if err := json.NewEncoder(&encountersBuffer).Encode(cacheEntry.Pokemon.Encounters); err != nil {
			continue
		}

		records = append(records, []string{
			id,
			name,
			typesBuffer.String(),
			statsBuffer.String(),
			encountersBuffer.String(),
			cacheEntry.CachedAt.Format(time.RFC3339),
		})
	}

	w.WriteAll(records)

	return w.Error()
}

func (c Cache) FindPokemon(input string) (*Pokemon, bool) {
	isID := true
	for _, char := range input {
		if char < '0' || char > '9' {
			isID = false
			break
		}
	}

	identifier := input
	if !isID {
		if ID, ok := c.NameMap[identifier]; ok {
			identifier = ID
		}
	}

	if entry, ok := c.IDMap[identifier]; ok {
		return &entry.Pokemon, true
	}
	return nil, false
}

func (c Cache) CachePokemon(pokemon Pokemon) error {
	pokemonID := strconv.Itoa(pokemon.ID)
	c.NameMap[pokemon.Name] = pokemonID
	c.IDMap[pokemonID] = CacheEntry{
		Pokemon:  pokemon,
		CachedAt: time.Now(),
	}

	return c.WriteToCSV()
}
