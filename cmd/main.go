package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":3000",
		db:   dbConfig{},
	}

	api := app{
		config: cfg,
	}

	// run the server
	if err := api.run(api.mount()); err != nil {
		slog.Error("failed to start the server", "error", err)
		os.Exit(1)
	}
}
