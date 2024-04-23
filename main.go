package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User struct represents a user model
type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

// Book struct represents a book model
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

var DB *gorm.DB
var secretKey = "your_secret_key"

func main() {
	// Connect to PostgreSQL database
	dsn := "host=localhost user=postgres password=your_password dbname=your_database_name port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db

	// Migrate the User and Book models
	DB.AutoMigrate(&User{}, &Book{})

	// Create a new Fiber instance
	app := fiber.New()

	// Authentication routes
	app.Post("/register", registerUser)
	app.Post("/login", loginUser)

	// Authenticated routes
	app.Use(authMiddleware)
	app.Get("/books", getAllBooks)
	app.Post("/books", createBook)
	app.Get("/books/:id", getBookByID)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}

// registerUser handler registers a new user
func registerUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	user.Password = string(hashedPassword)

	if err := DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(user)
}

// loginUser handler authenticates a user and generates a JWT token
func loginUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	var existingUser User
	if err := DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid username or password",
		})
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid username or password",
		})
	}

	// Generate a JWT token
	token, err := generateJWTToken(strconv.Itoa(int(existingUser.ID)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// authMiddleware middleware checks if the user is authenticated
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing Authorization header",
		})
	}

	token, err := validateJWTToken(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userID, err := strconv.Atoi(token.Claims.(jwt.MapClaims)["sub"].(string))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	c.Locals("userID", uint(userID))
	return c.Next()
}

// generateJWTToken generates a new JWT token
func generateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// validateJWTToken validates a JWT token
func validateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

// getAllBooks handler retrieves all books from the database
func getAllBooks(c *fiber.Ctx) error {
	var books []Book
	DB.Find(&books)
	return c.JSON(books)
}

// createBook handler creates a new book in the database
func createBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := DB.Create(book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(book)
}

// getBookByID handler retrieves a book by its ID from the database
func getBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}
	return c.JSON(book)
}

// updateBook handler updates a book in the database
func updateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}

	updatedBook := new(Book)
	if err := c.BodyParser(updatedBook); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Rating = updatedBook.Rating

	if err := DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(book)
}

// deleteBook handler deletes a book from the database
func deleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   fmt.Sprintf("Book with ID %s not found", id),
		})
	}

	if err := DB.Delete(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book with ID %s deleted successfully", id),
	})
}
