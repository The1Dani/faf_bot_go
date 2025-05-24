CREATE TABLE members (
    id INT DEFAULT nextval('id'::regclass) NOT NULL,
    chat_id BIGINT,
    member_id BIGINT,
    coefficient INT,
    pidor_coefficient INT,
    full_name VARCHAR(255),
    nick_name VARCHAR(255)
);
