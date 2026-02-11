package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TestProject_itk/internal/api"
	"TestProject_itk/internal/config"
	"TestProject_itk/internal/db"
	"TestProject_itk/internal/repository"
	"TestProject_itk/internal/service"
)

func main() {
	cfg, err := config.Load("config.env")
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	pool, err := db.NewPostgresPool(context.Background(), cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer pool.Close()

	log.Printf("Connected to database: %s\n", cfg.PostgresDSN())

	walletRepo := repository.NewWalletRepo(pool)
	walletSvc := service.NewWalletService(walletRepo)

	h := api.NewHandler(walletSvc)
	r := api.NewRouter(h)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("shutdown complete")
}
