package main

import (
	orderservice "gc1/order-service"
	productservice "gc1/product-service"
	userservice "gc1/user-service"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// Adjust the import path based on your module name
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load .env file from the root directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go orderservice.NewOrderService().StartCronJob()

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize Order Service
	orderService := orderservice.NewOrderService()
	orderService.RegisterRoutes(e)

	// Initialize Product Service
	productService := productservice.NewProductService()
	productService.RegisterRoutes(e)

	// Initialize User Service
	userService := userservice.NewUserService()
	userService.RegisterRoutes(e)

	// Start server
	log.Print("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
