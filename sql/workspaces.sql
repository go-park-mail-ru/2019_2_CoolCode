DROP TABLE IF EXISTS workspaces CASCADE;
CREATE TABLE workspaces
(
    ID        BIGSERIAL    NOT NULL
        PRIMARY KEY,
    name      VARCHAR(128) NULL,
    creatorID BIGINT,
    FOREIGN KEY (creatorID) REFERENCES users (ID)
);