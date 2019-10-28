DROP TABLE IF EXISTS chats_users CASCADE;
CREATE TABLE chats_users
(
    isAdmin BOOLEAN NULL ,
    chatID  BIGINT  NOT NULL,
    userID  BIGINT  NOT NULL,

    FOREIGN KEY (userID) REFERENCES users (ID),
    FOREIGN KEY (chatID) REFERENCES chats (ID),


);