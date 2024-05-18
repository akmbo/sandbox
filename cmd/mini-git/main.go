package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaolen/mini-git/internal/repository"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Subcommands: init, status")
		os.Exit(0)
	}

	switch os.Args[1] {

	case "init":
		initCmd.Parse(os.Args[2:])
		_, err := repository.Create(".")
		if err != nil {
			fmt.Printf("Error initializing new repository in current directory:\n\t%s\n", err)
			os.Exit(1)
		}
		fmt.Println("Intialized new repository")

	case "status":
		statusCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'status'")

	}
}
