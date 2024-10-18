-- +goose Up
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE    
);

INSERT INTO roles (name)
VALUES
    ('ROLE_USER'),
    ('ROLE_ADMIN');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,    
    role_id INTEGER NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
