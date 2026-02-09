package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/T-AKHMAD/gnotes-cli/internal/api"
)

func Notes(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "notes: missing subcommand")
		fmt.Fprintln(os.Stderr, "usage: gnotes notes list [--base-url ...]")
		return 1
	}

	switch args[0] {
	case "list":
		return notesList(args[1:])

	case "create":
		return notesCreate(args[1:])

	case "get":
		return notesGet(args[1:])

	case "delete":
		return notesDelete(args[1:])

	default:
		fmt.Fprintln(os.Stderr, "available: list, create, get, delete")
		return 1
	}
}

func notesList(args []string) int {
	fs := flag.NewFlagSet("notes list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes list:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.NotesList(ctx, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes list:", err)
		return 1
	}

	if len(res.Notes) == 0 {
		fmt.Println("no notes")
		return 0
	}

	for _, n := range res.Notes {
		fmt.Printf("%d | %s | %s\n", n.ID, n.CreatedAt, n.Title)
	}

	return 0
}

func notesCreate(args []string) int {
	fs := flag.NewFlagSet("notes create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")
	title := fs.String("title", "", "Note title")
	body := fs.String("body", "", "Note body")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *title == "" || *body == "" {
		fmt.Fprintln(os.Stderr, "notes create: --title and --body are required")
		fmt.Fprintln(os.Stderr, `example: gnotes notes create --title "t1" --body "b1"`)
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes create:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.CreateNote(ctx, token, api.CreateNoteRequest{
		Title: *title,
		Body:  *body,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes create:", err)
		return 1
	}

	fmt.Printf("id=%d\n", res.ID)
	return 0
}

func notesGet(args []string) int {
	fs := flag.NewFlagSet("notes get", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "notes get: <id> is required")
		fmt.Fprintln(os.Stderr, "example: gnotes notes get 2")
		return 1
	}

	id, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil || id <= 0 {
		fmt.Fprintln(os.Stderr, "notes get: invalid id")
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes get:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := client.GetNote(ctx, token, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes get:", err)
		return 1
	}

	fmt.Printf("id=%d\ncreated_at=%s\ntitle=%s\nbody=%s\n", n.ID, n.CreatedAt, n.Title, n.Body)
	return 0
}

func notesDelete(args []string) int {
	fs := flag.NewFlagSet("notes delete", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "notes delete: <id> is required")
		fmt.Fprintln(os.Stderr, "example: gnotes notes delete 2")
		return 1
	}

	id, err := strconv.ParseInt(rest[0], 10, 64)
	if err != nil || id <= 0 {
		fmt.Fprintln(os.Stderr, "notes delete: invalid id")
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "notes delete:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.DeleteNote(ctx, token, id); err != nil {
		fmt.Fprintln(os.Stderr, "notes delete:", err)
		return 1
	}

	fmt.Println("deleted")
	return 0
}
