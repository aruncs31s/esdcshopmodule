package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aruncs31s/esdcshopmodule/internal/application/usecase"
	"github.com/aruncs31s/esdcshopmodule/internal/domain/service"
	"github.com/aruncs31s/esdcshopmodule/internal/infrastructure/config"
	"github.com/aruncs31s/esdcshopmodule/internal/infrastructure/persistence"
	"github.com/aruncs31s/esdcshopmodule/internal/interface/http/handler"
	"github.com/aruncs31s/esdcshopmodule/internal/interface/http/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repositories (in-memory for now)
	productRepo := persistence.NewInMemoryProductRepository()
	orderRepo := persistence.NewInMemoryOrderRepository()
	cartRepo := persistence.NewInMemoryCartRepository()
	categoryRepo := persistence.NewInMemoryCategoryRepository()
	customerRepo := persistence.NewInMemoryCustomerRepository()

	// Initialize ID generator
	idGenerator := persistence.NewUUIDGenerator()

	// Initialize domain services
	orderService := service.NewOrderService(orderRepo, productRepo, cartRepo)

	// Initialize use cases
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, idGenerator)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productRepo, cartRepo, customerRepo, orderService, idGenerator)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo, idGenerator)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, idGenerator)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo, idGenerator)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	cartHandler := handler.NewCartHandler(cartUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	healthHandler := handler.NewHealthHandler(cfg.App.Version)

	// Setup router
	r := router.NewRouter(
		productHandler,
		orderHandler,
		cartHandler,
		categoryHandler,
		customerHandler,
		healthHandler,
	)

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r.Setup(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s server on %s", cfg.App.Name, addr)
		log.Printf("Environment: %s, Version: %s", cfg.App.Environment, cfg.App.Version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
