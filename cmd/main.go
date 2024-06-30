package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	salesHttp "go.dataflow.ru/service-sales/internal/adapters/http"
	"go.dataflow.ru/service-sales/internal/adapters/storage"
	"go.dataflow.ru/service-sales/internal/app/services"
	"go.dataflow.ru/service-sales/internal/config"
	"go.dataflow.ru/service-sales/pkg/logger"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := logger.NewLogger("debug")

	cfg, err := config.Read()
	if err != nil {
		logger.Panicf("cant read config: %v", err)
	}

	saleRepo := storage.New()
	saleService := services.NewSaleService(saleRepo)
	saleHandler := salesHttp.NewSalesHandler(saleService)
	srv, _ := NewServer(saleHandler)

	// Listen on a different Goroutine so the application doesn't stop here.
	if err := srv.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Panic(err)
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Perform application shutdown with a maximum timeout of 5 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		if err := srv.Shutdown(); err != nil {
			log.Fatalln(err)
		}
	}()

	select {
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			log.Fatalln("timeout exceeded, forcing shutdown")
		}

		os.Exit(0)
	}
}

func NewServer(h *salesHttp.SalesHandler) (_ *fiber.App, err error) {
	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Set up middlewares.

	server.Post("/data", h.AddSale)
	server.Get("/data", h.GetSales)

	return server, nil
}
