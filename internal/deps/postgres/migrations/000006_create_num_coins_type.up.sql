BEGIN;

-- Create a domain type to represent non-negative and NOT NULL amounts of coins
CREATE DOMAIN num_coins AS
    INT NOT NULL
    CONSTRAINT chk_non_negative_coins CHECK (VALUE >= 0);

-- Remove now unnecessary NOT NULL constraints
ALTER TABLE users
    ALTER COLUMN coins DROP NOT NULL;
ALTER TABLE coin_transactions
    ALTER COLUMN amount DROP NOT NULL;

-- Update coin columns to use the new type
ALTER TABLE users
    ALTER COLUMN coins TYPE num_coins;
ALTER TABLE coin_transactions
    ALTER COLUMN amount TYPE num_coins;

COMMIT;
