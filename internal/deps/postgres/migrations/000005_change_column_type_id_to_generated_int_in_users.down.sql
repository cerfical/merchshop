BEGIN;

ALTER TABLE users
    ALTER COLUMN id DROP IDENTITY;

-- Change the id column type to SERIAL
CREATE SEQUENCE users_id_seq AS INT;
ALTER TABLE users
    ALTER COLUMN id SET DEFAULT nextval('users_id_seq');
ALTER SEQUENCE users_id_seq OWNED BY users.id;

COMMIT;
