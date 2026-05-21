package main

import (
	"net/http"
	"time"
	"log/slog"

	_ "github.com/EYOB123695/ecom/docs"
	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
	"github.com/EYOB123695/ecom/internal/auth"
	"github.com/EYOB123695/ecom/internal/cart"
	authmw "github.com/EYOB123695/ecom/internal/middleware"
	"github.com/EYOB123695/ecom/internal/orders"
	"github.com/EYOB123695/ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	r.Get("/", health)

	queries := repo.New(app.db)

	productService := products.NewService(queries)
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.GetProductByID)
	r.Post("/products", productHandler.CreateProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	authService := auth.NewService(queries, app.config.jwtSecret)
	authHandler := auth.NewHandler(authService)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	cartService := cart.NewService(queries)
	cartHandler := cart.NewHandler(cartService)

	ordersService := orders.NewService(app.db, queries)
	ordersHandler := orders.NewHandler(ordersService)

	r.Group(func(r chi.Router) {
		r.Use(authmw.Auth(app.config.jwtSecret))
		r.Get("/users/me", authHandler.GetMe)

		// Cart routes
		r.Get("/cart", cartHandler.GetCart)
		r.Post("/cart", cartHandler.AddCartItem)
		r.Put("/cart/{product_id}", cartHandler.UpdateCartItem)
		r.Delete("/cart/{product_id}", cartHandler.DeleteCartItem)
		r.Delete("/cart", cartHandler.ClearCart)

		// Orders routes
		r.Post("/orders", ordersHandler.Checkout)
		r.Get("/orders", ordersHandler.ListOrders)
		r.Get("/orders/{id}", ordersHandler.GetOrder)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}
	app.logger.Info("Server has started", "addr", app.config.addr)
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Error("server failed to start", "error", err)
		return err
	}
	return nil
}

type application struct {
	config config
	db     *pgx.Conn
	logger *slog.Logger
}

type config struct {
	addr      string
	db        dbConfig
	jwtSecret string
}

type dbConfig struct {
	dsn string
}

// health godoc
// @Summary      Health check
// @Description  Returns a simple hello message
// @Tags         health
// @Produce      plain
// @Success      200 {string} string "hi"
// @Router       / [get]
func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
