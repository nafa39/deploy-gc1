package userservice

import (
	"context"
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

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type UserService struct {
	client *mongo.Client
}

func NewUserService() *UserService {
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

	return &UserService{client: client}
}

func (s *UserService) RegisterRoutes(e *echo.Echo) {
	// Get the users collection
	usersCollection := s.client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")

	// Routes
	e.POST("/users", s.createUser(usersCollection))
	e.GET("/users/:id", s.getUser(usersCollection))
	e.PUT("/users/:id", s.updateUser(usersCollection))
	e.DELETE("/users/:id", s.deleteUser(usersCollection))
	e.POST("/users/validate", s.validateToken)
}

// Create a new user
func (s *UserService) createUser(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		user.ID = primitive.NewObjectID()
		user.CreatedAt = time.Now()

		_, err := coll.InsertOne(context.Background(), user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		return c.JSON(http.StatusCreated, user)
	}
}

// Get a user by ID
func (s *UserService) getUser(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var user User
		err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Update a user by ID
func (s *UserService) updateUser(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var user User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		_, err = coll.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": user})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Delete a user by ID
func (s *UserService) deleteUser(coll *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// Validate token
func (s *UserService) validateToken(c echo.Context) error {
	type TokenRequest struct {
		Token string `json:"token"`
	}

	var req TokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Add your token validation logic here
	// For now, we'll just return a success response
	return c.JSON(http.StatusOK, map[string]bool{"valid": true})
}
