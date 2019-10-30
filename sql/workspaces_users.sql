DROP TABLE IF EXISTS workspaces_users CASCADE;
CREATE TABLE workspaces_users
(
    isAdmin     BOOLEAN NOT NULL,
    workspaceID BIGINT  NOT NULL,
    userID      BIGINT  NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (workspaceID) REFERENCES workspaces (ID) ON DELETE CASCADE
);