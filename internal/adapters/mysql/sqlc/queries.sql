-- name: GetProduct :one
SELECT
    *
FROM
    products
WHERE
    id = ?
LIMIT
    1;

-- name: ListProducts :many
SELECT
    *
FROM
    products;

-- name: CreateProduct :execresult
INSERT INTO
    products (name, price, quantity)
VALUES
    (?, ?, ?);

-- name: UpdateProductQty :execresult
UPDATE products
SET
    quantity = ?
WHERE
    id = ?;

-- name: DeleteProducts :exec
DELETE FROM products
WHERE
    id = ?;

-- name: GetOrder :one
SELECT
    *
FROM
    orders
WHERE
    id = ?
LIMIT
    1;

-- name: ListOrders :many
SELECT
    *
FROM
    orders;

-- name: CreateOrder :execresult
INSERT INTO
    orders (customer_id, total_price, status)
VALUES
    (?, ?, ?);

-- name: CreateOrderItems :execresult
INSERT INTO
    order_items (
        order_id,
        product_id,
        product_name,
        price,
        quantity
    )
VALUES
    (?, ?, ?, ?, ?);