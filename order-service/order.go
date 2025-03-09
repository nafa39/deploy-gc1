package orderservice

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type OrderService struct {
	client *mongo.Client
}

func NewOrderService() *OrderService {
	// Load .env file from the root directory
	// err := godotenv.Load("../.env") // Adjust the path if needed
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// MongoDB connection
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Verify MongoDB connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	return &OrderService{client: client}
}

func (s *OrderService) RegisterRoutes(e *echo.Echo) {
	// Routes
	e.POST("/orders", s.createOrder)
	e.GET("/orders/:id", s.getOrder)
	e.PUT("/orders/:id", s.updateOrder)
	e.DELETE("/orders/:id", s.deleteOrder)

	// Cron job to update order status daily
	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		collection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
		_, err := collection.UpdateMany(context.Background(), bson.M{"status": "Pending"}, bson.M{"$set": bson.M{"status": "Completed"}})
		if err != nil {
			log.Println("Error updating order statuses:", err)
		} else {
			log.Println("Order statuses updated")
		}
	})
	c.Start()
}

// Create a new order
func (s *OrderService) createOrder(c echo.Context) error {
	var order Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.Status = "Pending"

	// Insert order
	collection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order"})
	}

	return c.JSON(http.StatusCreated, order)
}

// Get an order by ID
func (s *OrderService) getOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Find order
	var order Order
	collection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

// Update an order by ID
func (s *OrderService) updateOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var order Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Update order
	collection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": order})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order"})
	}

	return c.JSON(http.StatusOK, order)
}

// Delete an order by ID
func (s *OrderService) deleteOrder(c echo.Context) error {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Delete order
	collection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("orders")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete order"})
	}

	return c.NoContent(http.StatusNoContent)
}
