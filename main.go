package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdFind = &cobra.Command{
		Use:   "find [id or name]",
		Short: "find a pokemen by id or name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pokemon := &Pokemon{PokeAPIBaseURL: "https://pokeapi.co/api/v2"}
			if err := pokemon.Find(args[0]); err != nil {
				fmt.Printf("failed to find pokemon, err: %v\n", err)
			}
			pokemonDetails, err := pokemon.PrettyPrint()
			if err != nil {
				fmt.Printf("failed to format pokemon details, displaying unstructured output\n%+v\n", pokemon)
			}

			fmt.Println(pokemonDetails)
		},
	}

	rootCmd := &cobra.Command{Use: "kidex"}
	rootCmd.AddCommand(cmdFind)
	rootCmd.Execute()
}
