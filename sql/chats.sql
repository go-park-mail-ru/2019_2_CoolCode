DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats
(
    ID            BIGSERIAL    NOT NULL PRIMARY KEY,
    isChannel     BOOLEAN      NOT NULL,
    totalMSGCount BIGINT       NOT NULL DEFAULT 0,
    name          VARCHAR(128) NULL,
    workspaceID   BIGINT       NULL,
    creatorID     BIGINT       NULL,
    FOREIGN KEY (workspaceID) REFERENCES workspaces (ID) ON DELETE CASCADE,
    FOREIGN KEY (creatorID) REFERENCES users (ID) ON DELETE SET NULL
);