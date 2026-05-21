-- name: ListProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, price_in_cents, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price_in_cents = $3, quantity = $4
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: UpdateProductQuantity :one
UPDATE products
SET quantity = quantity - $2
WHERE id = $1 AND quantity >= $2
RETURNING *;

-- name: GetCartItemsByUserID :many
SELECT c.*, p.name as product_name, p.price_in_cents as product_price, p.quantity as product_stock
FROM cart_items c
JOIN products p ON c.product_id = p.id
WHERE c.user_id = $1;

-- name: AddCartItem :one
INSERT INTO cart_items (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
RETURNING *;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $3
WHERE user_id = $1 AND product_id = $2
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE user_id = $1 AND product_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE user_id = $1;

-- name: CreateOrder :one
INSERT INTO orders (user_id, total, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 AND user_id = $2;

-- name: GetOrderItemsByOrderID :many
SELECT oi.*, p.name as product_name
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1;