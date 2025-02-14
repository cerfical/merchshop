CREATE TABLE coin_transactions(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    from_user_id INT NOT NULL REFERENCES users(id),
    to_user_id INT NOT NULL REFERENCES users(id),
    amount INT NOT NULL,

    CONSTRAINT chk_positive_amount CHECK (amount > 0),
    CONSTRAINT chk_different_users CHECK (from_user_id != to_user_id)
);
