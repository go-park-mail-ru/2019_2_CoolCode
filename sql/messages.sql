DROP TABLE IF EXISTS messages CASCADE;
CREATE TABLE "messages"(
    ID BIGSERIAL NOT NULL PRIMARY KEY ,
    type SMALLINT NOT NULL, --IN ('TEXT','PHOTO','VOICE')
    body TEXT NOT NULL,
    fileID BIGINT,
    chatID BIGINT NOT NULL,
    messageTime TIME,
    authorID BIGINT NOT NULL,
    hideForAuthor bool default false,
    FOREIGN KEY (authorID)  REFERENCES users(ID),
    FOREIGN KEY (chatID)  REFERENCES chats(ID)
)