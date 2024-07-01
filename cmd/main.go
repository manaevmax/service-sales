package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.dataflow.ru/service-sales/pkg/logger"

	salesHttp "go.dataflow.ru/service-sales/internal/adapters/http"
	"go.dataflow.ru/service-sales/internal/adapters/storage"
	"go.dataflow.ru/service-sales/internal/app/services"
	"go.dataflow.ru/service-sales/internal/config"
)

const (
	readTimeout  = 1 * time.Second
	writeTimeout = 1 * time.Second
)

func main() {
	logger := logger.NewLogger("debug")

	cfg, err := config.Read()
	if err != nil {
		logger.Panicf("cant read config: %v", err)
	}

	saleRepo := storage.New(logger)
	saleService := services.NewSaleService(saleRepo, logger)
	saleHandler := salesHttp.New(saleService)
	srv := NewServer(saleHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c

		logger.Info("Gracefully shutting down...")

		_ = srv.Shutdown()
	}()

	if err := srv.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Panic(err)
	}
}

func NewServer(h *salesHttp.SalesHandler) *fiber.App {
	server := fiber.New(fiber.Config{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	server.Post("/data", h.AddSale)
	server.Get("/data", h.GetSales)
	server.Post("/calculate", h.CalculateTotalSum)

	return server
}
