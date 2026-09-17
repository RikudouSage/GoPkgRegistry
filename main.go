package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/cors"
	"go.chrastecky.dev/go-pkg-repository/cfg"
	"go.chrastecky.dev/go-pkg-repository/db"
	"go.chrastecky.dev/go-pkg-repository/handlers"
	appMiddleware "go.chrastecky.dev/go-pkg-repository/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/lo"
)

var globalConfig *cfg.GlobalConfig
var database *db.Client

func getRouter() chi.Router {
	authMiddleware := appMiddleware.RequiresAuthorizationMiddleware(globalConfig.AdminAPIKey)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.GetHead)
	if globalConfig.FrontendURL != "" {
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{globalConfig.FrontendURL},
			AllowedHeaders: []string{"Authorization", "Content-Type"},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		}))
	}

	router.With(authMiddleware).Route("/admin", func(router chi.Router) {
		router.Get("/packages", func(writer http.ResponseWriter, request *http.Request) {
			handlers.GetPackagesHandler(writer, request, database)
		})
		router.Get("/packages/by-id/{id}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.GetPackageByIDHandler(writer, request, database)
		})
		router.Get("/packages/*", func(writer http.ResponseWriter, request *http.Request) {
			handlers.GetPackageHandler(writer, request, database)
		})
		router.Post("/packages", func(writer http.ResponseWriter, request *http.Request) {
			handlers.StorePackageHandler(writer, request, database)
		})
	})

	router.Get("/*", func(writer http.ResponseWriter, req *http.Request) {
		handlers.PackageHandler(writer, req, globalConfig, database)
	})

	return router
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	server := &http.Server{
		Addr:              ":" + fmt.Sprint(globalConfig.Port),
		Handler:           getRouter(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	go func() {
		log.Println("Starting server on " + fmt.Sprint(globalConfig.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-gracefulShutdown
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Println("Server gracefully stopped")
}

func init() {
	globalConfig = lo.Must(cfg.GetGlobalConfig())
	database = lo.Must(db.NewClient(globalConfig))
}
