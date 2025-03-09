package productservice

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Price     float64            `json:"price" bson:"price"`
	Stock     int                `json:"stock" bson:"stock"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type ProductService struct {
	client *mongo.Client
}

func NewProductService() *ProductService {
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
		log.Fatal(err)
	}

	// Verify MongoDB connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	return &ProductService{client: client}
}

func (s *ProductService) RegisterRoutes(e *echo.Echo) {
	// Get the products collection
	productsCollection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("products")

	// Routes
	e.POST("/products", s.createProduct(productsCollection))
	e.GET("/products/:id", s.getProduct(productsCollection))
	e.PUT("/products/:id", s.updateProduct(productsCollection))
	e.DELETE("/products/:id", s.deleteProduct(productsCollection))
}

// Create a new product
func (s *ProductService) createProduct(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var product Product
		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		product.ID = primitive.NewObjectID()
		product.CreatedAt = time.Now()

		// Start a session for transaction (if needed)
		session, err := coll.Database().Client().StartSession()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start session"})
		}
		defer session.EndSession(context.Background())

		// Transaction callback
		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			_, err := coll.InsertOne(sessCtx, product)
			if err != nil {
				return nil, err
			}
			return product, nil
		}

		// Execute transaction
		result, err := session.WithTransaction(context.Background(), callback)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product"})
		}

		return c.JSON(http.StatusCreated, result)
	}
}

// Get a product by ID
func (s *ProductService) getProduct(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var product Product
		err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch product"})
		}

		return c.JSON(http.StatusOK, product)
	}
}

// Update a product by ID
func (s *ProductService) updateProduct(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var product Product
		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		// Start a session for transaction (if needed)
		session, err := coll.Database().Client().StartSession()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start session"})
		}
		defer session.EndSession(context.Background())

		// Transaction callback
		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			_, err := coll.UpdateOne(sessCtx, bson.M{"_id": objID}, bson.M{"$set": product})
			if err != nil {
				return nil, err
			}
			return product, nil
		}

		// Execute transaction
		result, err := session.WithTransaction(context.Background(), callback)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
		}

		return c.JSON(http.StatusOK, result)
	}
}

// Delete a product by ID
func (s *ProductService) deleteProduct(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		// Start a session for transaction (if needed)
		session, err := coll.Database().Client().StartSession()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start session"})
		}
		defer session.EndSession(context.Background())

		// Transaction callback
		callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
			_, err := coll.DeleteOne(sessCtx, bson.M{"_id": objID})
			if err != nil {
				return nil, err
			}
			return nil, nil
		}

		// Execute transaction
		_, err = session.WithTransaction(context.Background(), callback)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
