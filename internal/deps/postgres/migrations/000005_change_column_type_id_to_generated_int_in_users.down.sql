BEGIN;

ALTER TABLE users
ALTER COLUMN id DROP IDENTITY;

CREATE SEQUENCE users_id_seq AS INT;

-- Change the id column type to SERIAL
ALTER TABLE users
ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

ALTER SEQUENCE users_id_seq OWNED BY users.id;

COMMIT;
