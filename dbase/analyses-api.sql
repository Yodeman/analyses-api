CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "files" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "data" text NOT NULL,
  "changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "files" ("username");

ALTER TABLE "users" ADD CONSTRAINT "user_email_constraint" UNIQUE ("username", "email");

ALTER TABLE "files" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
