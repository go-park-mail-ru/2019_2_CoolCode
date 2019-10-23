DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    ID           BIGSERIAL    NOT NULL
        PRIMARY KEY,
    username     VARCHAR(32)  NOT NULL UNIQUE,
    email        VARCHAR(128) NOT NULL UNIQUE,
    name         VARCHAR(128),
    status       VARCHAR(32),
    phone        VARCHAR(12),
    passwordHash BYTEA        NOT NULL,
    salt         BYTEA        NOT NULL
);