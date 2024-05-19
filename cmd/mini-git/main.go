package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaolen/mini-git/internal/objects"
	"github.com/aaolen/mini-git/internal/repository"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	writeStringCmd := flag.NewFlagSet("write-string", flag.ExitOnError)
	readObjectCmd := flag.NewFlagSet("read-object", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Subcommands: init, write-string, read-object")
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

	case "write-string":
		writeStringCmd.Parse(os.Args[2:])
		input := writeStringCmd.Arg(0)
		if input == "" {
			fmt.Println("Expected argument: input")
			os.Exit(1)
		}

		r, err := repository.Discover(".")
		if err != nil {
			fmt.Println("Not inside repository")
			os.Exit(1)
		}
		checksum, err := objects.WriteBlob(r, input)
		if err != nil {
			fmt.Printf("Error writing object:\n\t%s\n", err)
			os.Exit(1)
		}
		fmt.Println(checksum)

	case "read-object":
		readObjectCmd.Parse(os.Args[2:])
		checksum := readObjectCmd.Arg(0)
		if checksum == "" {
			fmt.Println("Expected argument: checksum")
			os.Exit(1)
		}

		r, err := repository.Discover(".")
		if err != nil {
			fmt.Println("Not inside repository")
			os.Exit(1)
		}

		output, err := objects.ReadBlob(r, checksum)
		if err != nil {
			fmt.Printf("Error retrieving object:\n\t%s\n", err)
			os.Exit(1)
		}
		fmt.Println(output)

	}

}
