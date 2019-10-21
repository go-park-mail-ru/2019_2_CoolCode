DROP TABLE IF EXISTS "chats_users" CASCADE;
CREATE TABLE "chats_users"
(
    isAdmin     BOOLEAN NOT NULL,
    workspaceID BIGINT  NOT NULL,
    userID      BIGINT  NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (ID),
    FOREIGN KEY (workspaceID) REFERENCES workspaces (ID)
);