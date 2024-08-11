package routes

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Trisamudrisvara/goTrip/db"
)

// getTrips retrieves all trips from the database
func (r *Repo) ListTrips(c *fiber.Ctx) error {
	trips, err := r.Queries.ListTrips(r.Ctx)

	if err != nil {
		log.Println("Error in getting trips in GetTrips db function:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	if len(trips) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{"error": "no trips found"})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": &trips,
	})
}

// getTrip retrieves a single trip by ID
func (r *Repo) getTrip(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid UUID") {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidID)
		}

		log.Println("Error in parsing uuid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	id := pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}

	trip, err := r.Queries.GetTrip(r.Ctx, id)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidID)
		}

		log.Println("Error in getting trip in GetTrip db function:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": &trip,
	})
}

// createTrip adds a new trip to the database
func (r *Repo) createTrip(c *fiber.Ctx) error {
	// Extract trip details from form data
	name := c.FormValue("name")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")
	destinationId := c.FormValue("destination_id")

	// Validate input
	if name == "" || startDate == "" || endDate == "" || destinationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiberUndefinedParamError)
	}

	// Parse destination UUID
	Uuid, err := uuid.Parse(destinationId)

	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid UUID") {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidDestinationID)
		}

		log.Println("Error in parsing uuid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Prepare trip data for database insertion
	trip := db.CreateTripParams{
		ID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		DestinationID: pgtype.UUID{
			Bytes: Uuid,
			Valid: true,
		},
	}

	// Create trip in database
	err = r.Queries.CreateTrip(r.Ctx, trip)

	if err != nil {
		if err.Error() == "ERROR: insert or update on table \"trip\" violates foreign key constraint \"trip_destination_id_fkey\" (SQLSTATE 23503)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidDestinationID)
		}

		log.Println("Error in creating trip in CreateTrip db function:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "trip has been added"})
}

// updateTrip modifies an existing trip in the database
func (r *Repo) updateTrip(c *fiber.Ctx) error {
	// Extract trip details from form data
	id := c.FormValue("id")
	name := c.FormValue("name")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")
	destinationId := c.FormValue("destination_id")

	// Validate input
	if id == "" || name == "" || startDate == "" || endDate == "" || destinationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiberUndefinedParamError)
	}

	// Parse destination UUID
	destinationUuid, err := uuid.Parse(destinationId)

	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid UUID") {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidDestinationID)
		}

		log.Println("Error in parsing uuid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Parse trip UUID
	uuid, err := uuid.Parse(id)

	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid UUID") {
			return c.Status(fiber.StatusBadRequest).JSON(
				&fiber.Map{"error": "invalid trip id"})
		}

		log.Println("Error in parsing uuid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	// Prepare trip data for database update
	trip := db.UpdateTripParams{
		ID: pgtype.UUID{
			Bytes: uuid,
			Valid: true,
		},
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		DestinationID: pgtype.UUID{
			Bytes: destinationUuid,
			Valid: true,
		},
	}

	// Update trip in database
	err = r.Queries.UpdateTrip(r.Ctx, trip)

	if err != nil {
		if err.Error() == "ERROR: insert or update on table \"trip\" violates foreign key constraint \"trip_destination_id_fkey\" (SQLSTATE 23503)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidDestinationID)
		}

		log.Println("Error in updating trip in UpdateTrip db function:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "trip has been updated"})
}

// deleteTrip removes a trip from the database by ID
func (r *Repo) deleteTrip(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("id"))

	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid UUID") {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidID)
		}

		log.Println("Error in parsing uuid:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	id := pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}

	err = r.Queries.DeleteTrip(r.Ctx, id)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.Status(fiber.StatusBadRequest).JSON(fiberInvalidID)
		}

		log.Println("Error in deleting trip in DeleteTrip db function:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiberUnknownError)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "trip has been deleted"})
}
