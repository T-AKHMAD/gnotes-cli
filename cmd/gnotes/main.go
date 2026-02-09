package main

import (
	"fmt"
	"os"

	"github.com/T-AKHMAD/gnotes-cli/internal/cli"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) < 2 {
		printHelp()
		return 1
	}
	switch args[1] {

	case "help", "-h", "--help":
		printHelp()
		return 0

	case "version", "-v", "--version":
		fmt.Println(version)
		return 0

	case "login":
		return cli.Login(args[2:])

	case "me":
		return cli.Me(args[2:])

	case "notes":
		return cli.Notes(args[2:])

	case "logout":
		return cli.Logout(args[2:])

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[1])
		return 1
	}
}

func printHelp() {
	fmt.Println(`gnotes - CLI client for gopher-notes

	Usage:
	gnotes <command> [args]

	Commands:
	help       Show help
	version    Show version

	Examples:
	gnotes version
	`)
}
