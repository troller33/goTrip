package routes

import (
	"fmt"
	"log"
	"time"
	"trip/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// getCsrf retrieves the CSRF token from the context and returns it
func getCsrfToken(c *fiber.Ctx) error {
	csrfToken, ok := c.Locals("csrf").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{"error": "Failed to get csrf token"})
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Map{"csrf": csrfToken})
}

// login handles user authentication and JWT token generation
func (r *Repo) login(c *fiber.Ctx) error {
	// Extract email and password from form data
	email := c.FormValue("email")
	pass := c.FormValue("password")

	// Validate input
	if pass == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiberUndefinedParamError)
	}

	// Retrieve user's password hash from database
	GetPass, err := r.Queries.GetPass(r.Ctx, email)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiberInvalidEmailPass)
		}

		log.Println("error in getting user info from db in GetPass function:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Construct user UUID and append to password for verification
	usrUuid := fmt.Sprintf("%x-%x-%x-%x-%x", GetPass.ID.Bytes[0:4], GetPass.ID.Bytes[4:6], GetPass.ID.Bytes[6:8], GetPass.ID.Bytes[8:10], GetPass.ID.Bytes[10:16])
	pass += usrUuid
	err = bcrypt.CompareHashAndPassword([]byte(GetPass.Password), []byte(pass))

	// Check if password is correct
	if err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiberInvalidEmailPass)
		}

		log.Println("error in comparing hash password:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Check if user is owner
	isOwner := usrUuid == ownerUuid

	// Create JWT claims
	claims := jwt.MapClaims{
		"email": email,
		"name":  GetPass.Name,
		"admin": GetPass.Admin,
		"owner": isOwner,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create and sign JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	jwtToken, err := token.SignedString(secret)

	if err != nil {
		log.Println("error in signing JWT key:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.JSON(fiber.Map{"jwt": jwtToken})
}

// register handles user registration
func (r *Repo) register(c *fiber.Ctx) error {
	// Extract registration details from form data
	name := c.FormValue("name")
	pass := c.FormValue("password")
	email := c.FormValue("email")

	// Validate input
	if name == "" || pass == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiberUndefinedParamError)
	}

	// Generate UUID and hash password
	uuid := uuid.New()
	pass += uuid.String()
	password, err := bcrypt.GenerateFromPassword([]byte(pass), 12)

	if err != nil {
		log.Println("Error hashing password:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Prepare user data for database insertion
	usr := db.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		},
		Email:    email,
		Name:     name,
		Password: string(password),
	}

	// Create user in database
	err = r.Queries.CreateUser(r.Ctx, usr)

	if err != nil {
		log.Println("Error in creating new user in CreateUser db function:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "user has been added"})
}
