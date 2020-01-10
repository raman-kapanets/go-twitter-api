CREATE TABLE tweets
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    message TEXT
);