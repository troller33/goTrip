package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Trisamudrisvara/goTrip/db"
)

// aboutUser retrieves and returns user information from the JWT token
func aboutUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	email := claims["email"].(string)
	return c.SendString("Welcome " + name + "\n" + email)
}

// updateUser handles the user update process
func (r *Repo) updateUser(c *fiber.Ctx) error {
	// Extract user information from the JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	// Get form values for user update
	oldEmail := c.FormValue("old_email")
	newEmail := c.FormValue("new_email")
	name := c.FormValue("name")

	// Validate input
	if oldEmail == "" || newEmail == "" || name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiberUndefinedParamError)
	}

	// Check if the user is authorized to make this update
	if email != oldEmail {
		return c.Status(fiber.StatusUnauthorized).JSON(fiberUnauthorizedError)
	}

	// Prepare user data for database update
	usr := db.UpdateUserParams{
		Email:   oldEmail,
		Email_2: newEmail,
		Name:    name,
	}

	// Update user in database
	err := r.Queries.UpdateUser(r.Ctx, usr)

	if err != nil {
		log.Println("Error in updating user in UpdateUser db function:", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiberUndefinedParamError)
	}

	// Create new JWT claims with updated information
	claims = jwt.MapClaims{
		"email": newEmail,
		"name":  name,
		"admin": claims["admin"],
		"owner": claims["owner"],
		"exp":   claims["exp"],
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
