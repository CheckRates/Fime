CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "createdAt" timestamptz DEFAULT (now())
);

CREATE TABLE "images" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "url" varchar NOT NULL,
  "owner" bigserial
);

CREATE TABLE "tags" (
  "id" bigserial PRIMARY KEY,
  "tag" varchar(30) UNIQUE NOT NULL
);

CREATE TABLE "image_tags" (
  "image_id" bigint REFERENCES "images" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  "tag_id" bigint REFERENCES "tags" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  PRIMARY KEY (image_id, tag_id)
);

ALTER TABLE "images" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

CREATE INDEX ON "users" ("name");

CREATE INDEX ON "images" ("owner");

CREATE INDEX ON "images" ("name");
