-- name: CreateTransaction :exec
INSERT INTO transactions (account_id, amount, type, transaction_at) VALUES ($1, $2, $3, $4);

-- name: GetTransactionsByAccountID :many
SELECT id, account_id, amount, type, transaction_at FROM transactions WHERE account_id = $1;

-- name: CreateTransactions :exec
INSERT INTO transactions (account_id, amount, type, transaction_at) VALUES (
    unnest($1::int[]),
    unnest($2::float8[]),
    unnest($3::varchar[]),
    unnest($4::date[])
);