-- name: CreateOrder :one
INSERT INTO orders (code, customer_code, total_value, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product, quantity, price, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderByCode :one
SELECT * FROM orders
WHERE code = $1;

-- name: GetOrdersByCustomerCode :many
SELECT * FROM orders
WHERE customer_code = $1
ORDER BY created_at DESC;

-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: GetTotalByOrderCode :one
SELECT total_value FROM orders
WHERE code = $1;

-- name: CountOrdersByCustomer :one
SELECT COUNT(*) FROM orders
WHERE customer_code = $1;
