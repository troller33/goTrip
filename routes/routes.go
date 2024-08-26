package routes

import (
	"context"
	"os"

	"github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"github.com/Trisamudrisvara/goTrip/db"
)

type CheckIfOwner bool

type Repo struct {
	Ctx     context.Context
	Queries *db.Queries
}

var (
	secret    []byte
	ownerUuid string

	// Defining Errors
	fiberUnknownError         = &fiber.Map{"error": "some unknown error occured"}
	fiberUndefinedParamError  = &fiber.Map{"error": "some params are undefined"}
	fiberUnauthorizedError    = &fiber.Map{"error": "unauthorized"}
	fiberInvalidID            = &fiber.Map{"error": "invalid id"}
	fiberInvalidDestinationID = &fiber.Map{"error": "invalid destination id"}
	fiberInvalidEmailPass     = &fiber.Map{"error": "invalid email or password"}
)

func (r *Repo) SetupRoutes(app *fiber.App) {
	// loads neccessary environment variables
	loadEnvVars()

	// Prometheus
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Auth
	login := app.Group("/login")
	login.Get("", getCsrfToken)
	login.Post("", r.login)
	app.Post("/register", r.register)

	// For testing csrf
	app.Post("", hello)

	// initializing /destination route
	destination := app.Group("/destination")
	destination.Get("", r.ListDestinations)
	destination.Get("/:id", r.getDestination)

	// initializing /trip route
	trip := app.Group("/trip")
	trip.Get("", r.ListTrips)
	trip.Get("/:id", r.getTrip)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: secret},
	}))

	// JWT Routes below

	// /user route
	usr := app.Group("/user")
	// usr.Get("", aboutUser) // GET isn't protected by CSRF
	usr.Post("", aboutUser)
	usr.Put("", r.updateUser)

	// only owner can promote user to admin
	// or demote admin to user with additional admin=demote in form
	// app.Get("admin/:email", r.promoteAdmin) // GET isn't protected by CSRF
	checkOwner := CheckIfOwner(true)
	app.Post("/admin", checkOwner.checkIfAdmin, r.promoteAdmin)

	// checks for admin instead of owner
	checkOwner = false
	// Middleware to check if user is admin
	app.Use(checkOwner.checkIfAdmin)
	// test if middleware is working properly
	app.Get("", hello)

	// changes destination
	destination.Post("", r.createDestination)
	destination.Put("", r.updateDestination)
	destination.Delete("/:id", r.deleteDestination)

	// changes trip
	trip.Post("", r.createTrip)
	trip.Put("", r.updateTrip)
	trip.Delete("/:id", r.deleteTrip)
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

// loads neccessary environment variables
func loadEnvVars() {
	// Get the secret signing Key for JWT
	secret = []byte(os.Getenv("SECRET"))
	// Get uuid of owner
	ownerUuid = os.Getenv("OWNER_UUID")

	// log.Println(secret)
	// log.Println(ownerUuid)
}
