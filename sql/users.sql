DROP TABLE IF EXISTS "users" CASCADE;
CREATE TABLE "users"
(
    id BIGSERIAL NOT NULL
            PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(128) NOT NULL,
    name VARCHAR(128),
    password VARCHAR(128) NOT NULL,
    status VARCHAR(32),
    phone VARCHAR(12)
);