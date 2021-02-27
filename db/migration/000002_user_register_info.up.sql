ALTER TABLE users ADD COLUMN email varchar UNIQUE NOT NULL;
ALTER TABLE users ADD COLUMN hashedPassword varchar NOT NULL;
ALTER TABLE users ADD COLUMN passwordChangedAt timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z';