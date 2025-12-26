-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    total_price INT NOT NULL CHECK (total_price >= 0),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE TABLE IF EXIST orders
-- +goose StatementEnd
