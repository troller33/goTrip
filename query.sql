-- name: GetPass :one
SELECT id, password, name, admin FROM users
 WHERE email = $1 LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
  id, email, name, password
) VALUES (
  $1, $2, $3, $4
);

-- name: UpdateUser :exec
UPDATE users
 SET email = $2,
 name = $3
WHERE email = $1;


-- name: PromoteAdmin :exec
UPDATE users SET admin = true
 WHERE email = $1;

-- name: DemoteAdmin :exec
UPDATE users SET admin = false
 WHERE email = $1;

-- name: CreateDestination :exec
INSERT INTO destination (
  id, name, description, attraction
) VALUES (
  $1, $2, $3, $4
);

-- name: ListDestinations :many
SELECT * FROM destination;

-- name: GetDestination :one
SELECT name, description, attraction FROM destination
 WHERE id = $1 LIMIT 1;

-- name: UpdateDestination :exec
UPDATE destination
 SET name = $2,
 description = $3,
 attraction = $4
WHERE id = $1;

-- name: DeleteDestination :exec
DELETE FROM destination
 WHERE id = $1;


-- name: CreateTrip :exec
INSERT INTO trip (
  id, name, start_date, end_date, destination_id
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: ListTrips :many
SELECT * FROM trip;

-- name: GetTrip :one
SELECT name, start_date, end_date, destination_id FROM trip
 WHERE id = $1 LIMIT 1;

-- name: UpdateTrip :exec
UPDATE trip
 SET name = $2,
 start_date = $3,
 end_date = $4,
 destination_id = $5
WHERE id = $1;

-- name: DeleteTrip :exec
DELETE FROM trip
 WHERE id = $1;
