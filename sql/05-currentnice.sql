CREATE TABLE currentnice (
    id INT DEFAULT nextval('id'::regclass) NOT NULL,
    chat_id BIGINT,
    member_id BIGINT,
    timestamp BIGINT
);
