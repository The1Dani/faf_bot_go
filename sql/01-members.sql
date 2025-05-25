CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT,
    member_id BIGINT,
    coefficient INT,
    pidor_coefficient INT,
    full_name VARCHAR(255),
    nick_name VARCHAR(255)
);
