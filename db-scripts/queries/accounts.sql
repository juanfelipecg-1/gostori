-- name: CreateAccount :one
INSERT INTO accounts (email) VALUES ($1) RETURNING id;

-- name: GetAccountByID :one
SELECT id, email FROM accounts WHERE id = $1;

-- name: GetAccountByEmail :one
SELECT id, email FROM accounts WHERE email = $1;