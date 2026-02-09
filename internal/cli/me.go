package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/T-AKHMAD/gnotes-cli/internal/api"
)

func Me(args []string) int {
	fs := flag.NewFlagSet("me", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	baseURL := fs.String("base-url", "http://localhost:8080", "API base URL")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	token, err := LoadToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "me:", err)
		fmt.Fprintln(os.Stderr, "tip: run `gnotes login --email ... --password ...` first")
		return 1
	}

	client := api.NewClient(*baseURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	me, err := client.Me(ctx, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "me:", err)
		return 1
	}

	fmt.Printf("id=%d\nemail=%s\n", me.ID, me.Email)
	return 0
}
