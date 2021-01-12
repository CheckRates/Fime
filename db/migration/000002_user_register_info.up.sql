ALTER TABLE users ADD COLUMN email varchar NOT NULL;
ALTER TABLE users ADD COLUMN password varchar NOT NULL;
ALTER TABLE users ADD CONSTRAINT minPasswordLen CHECK (length(password) >= 8);