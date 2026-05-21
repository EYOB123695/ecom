// @title           Ecom API
// @version         1.0
// @description     Ecommerce backend API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/EYOB123695/ecom/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("DB_DSN", "host=localhost port=5433 user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
		jwtSecret: env.GetString("JWT_SECRET", "dev-secret-change-in-production"),
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	api := application{
		config: cfg,
		db:     conn,
		logger: logger,
	}
	logger.Info("database connected successfully")

	if err := api.run(api.mount()); err != nil {
		logger.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
