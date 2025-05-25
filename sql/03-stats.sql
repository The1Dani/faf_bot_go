CREATE TABLE stats (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT,
    member_id BIGINT,
    count INT
);
