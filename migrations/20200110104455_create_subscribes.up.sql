CREATE TABLE subscribes
(
    id            BIGSERIAL NOT NULL PRIMARY KEY,
    subscriber    BIGINT REFERENCES users (id) ON DELETE CASCADE,
    subscribed_to BIGINT REFERENCES users (id) ON DELETE CASCADE,
    unique (subscriber, subscribed_to)
);