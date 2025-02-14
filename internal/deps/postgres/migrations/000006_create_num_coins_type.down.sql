BEGIN;

-- Restore the original column types
ALTER TABLE users
ALTER COLUMN coins TYPE INT;

ALTER TABLE coin_transactions
ALTER COLUMN amount TYPE INT;

-- Restore the NOT NULL constraints
ALTER TABLE users
ALTER COLUMN coins SET NOT NULL;

ALTER TABLE coin_transactions
ALTER COLUMN amount SET NOT NULL; 

DROP DOMAIN num_coins;

COMMIT;
