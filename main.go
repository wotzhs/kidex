package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cache := Cache{}
	if err := cache.Restore(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("failed to restore cache err: %v\n", err)
			fmt.Println("Re-creating a new cache")
		}
	}

	var cmdFind = &cobra.Command{
		Use:   "find [id or name]",
		Short: "find a pokemen by id or name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pokemon, found := cache.FindPokemon(args[0])

			if !found {
				pokemon = &Pokemon{PokeAPIBaseURL: "https://pokeapi.co/api/v2"}
				fmt.Println("fetching from " + pokemon.PokeAPIBaseURL)

				if err := pokemon.Find(args[0]); err != nil {
					fmt.Printf("failed to find pokemon, err: %v\n", err)
					return
				}
				if err := cache.CachePokemon(*pokemon); err != nil {
					fmt.Println("failed to update cache")
				}
			}

			if found {
				fmt.Println("served from cache")
			}

			pokemonDetails, err := pokemon.PrettyPrint()
			if err != nil {
				fmt.Printf("failed to format pokemon details, displaying unstructured output\n%+v\n", pokemon)
				return
			}

			fmt.Println(pokemonDetails)
		},
	}

	rootCmd := &cobra.Command{Use: "kidex"}
	rootCmd.AddCommand(cmdFind)
	rootCmd.Execute()
}
