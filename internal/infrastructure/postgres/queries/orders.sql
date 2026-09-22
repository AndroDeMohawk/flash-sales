-- internal/infrastructure/postgres/queries/orders.sql

-- name: CreateOrder :one
INSERT INTO orders (id, user_id, ticket_id, quantity, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;