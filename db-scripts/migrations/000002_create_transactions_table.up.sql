CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    amount DOUBLE PRECISION NOT NULL,
    type VARCHAR(10) NOT NULL,
    transaction_at DATE NOT NULL
);
