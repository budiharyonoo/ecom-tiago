-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    quantity INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE TABLE IF EXISTS products
-- +goose StatementEnd
