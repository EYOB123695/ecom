package main

import (
	"log"
	"log/slog"
	"os"
)


func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: "",
		},
	}

	api := application{
		config: cfg,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)


	if err := api.run(api.mount()); err != nil {
		log.Println("Server has failed to start")
		os.Exit(1)
	} 
}
