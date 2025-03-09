package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/joho/godotenv/autoload"
)

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	ProductID string             `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Total     float64            `json:"total" bson:"total"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Verify MongoDB connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/orders", createOrder)
	e.GET("/orders/:id", getOrder)
	e.PUT("/orders/:id", updateOrder)
	e.DELETE("/orders/:id", deleteOrder)

	// Cron job to update order status daily
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
		_, err := collection.UpdateMany(ctx, bson.M{"status": "Pending"}, bson.M{"$set": bson.M{"status": "Completed"}})
		if err != nil {
			log.Println("Error updating order statuses:", err)
		} else {
			log.Println("Order statuses updated")
		}
	})
	c.Start()

	// Start server
	e.Logger.Fatal(e.Start(":3003"))
}

// Create a new order
func createOrder(c echo.Context) error {
	var order Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.Status = "Pending"

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to MongoDB"})
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Insert order
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err = collection.InsertOne(ctx, order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}

	return c.JSON(http.StatusCreated, order)
}

// Get an order by ID
func getOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to MongoDB"})
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Find order
	var order Order
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

// Update an order by ID
func updateOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var order Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to MongoDB"})
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Update order
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": order})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order"})
	}

	return c.JSON(http.StatusOK, order)
}

// Delete an order by ID
func deleteOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to MongoDB"})
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Delete order
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete order"})
	}

	return c.NoContent(http.StatusNoContent)
}
