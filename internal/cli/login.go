package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/T-AKHMAD/gnotes-cli/internal/api"
)

func Login(args []string) int {

	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")
	email := fs.String("email", "", "User email")
	password := fs.String("password", "", "User password")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *email == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "login: --email and --password are required")
		fmt.Fprintln(os.Stderr, "example: gnotes login --email a@b.com --password 123")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.Login(ctx, api.LoginRequest{
		Email:    *email,
		Password: *password,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "login:", err)
		return 1
	}

	if err := SaveToken(res.Token); err != nil {
		fmt.Fprintln(os.Stderr, "login: failed to save token:", err)
		return 1
	}

	fmt.Fprintln(os.Stdout, "login: ok")
	fmt.Printf("token saved\nexpires_at=%s\n", res.ExpiresAt)
	return 0
}
