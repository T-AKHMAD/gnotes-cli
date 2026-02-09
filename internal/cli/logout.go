package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/T-AKHMAD/gnotes-cli/internal/api"
)

func Logout(args []string) int {
	fs := flag.NewFlagSet("logout", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "logout:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Logout(ctx, token); err != nil {
		fmt.Fprintln(os.Stderr, "logout:", err)
		return 1
	}

	p, err := tokenPath()
	if err == nil {
		_ = os.Remove(p)
	}

	fmt.Println("logged out")
	return 0
}
