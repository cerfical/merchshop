-- Create a table to store merch items purchased by each user
CREATE TABLE user_inventories(
    user_id INT NOT NULL REFERENCES users(id),
    merch TEXT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, merch),
    CONSTRAINT chk_positive_quantity CHECK (quantity >= 0)
);
