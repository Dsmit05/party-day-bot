CREATE TABLE users
(
    id         serial PRIMARY KEY,
    tg_id      bigint,
    chat_id    bigint,
    role       text NOT NULL DEFAULT 'Ordinary',
    first_name text,
    last_name  text,
    user_name  text
);

create unique index tg_id_index
    on users (tg_id);

CREATE TABLE files
(
    id      serial PRIMARY KEY,
    tg_id   text,
    url     text,
    user_tg_id bigint,
    created_at    timestamp with time zone NOT NULL DEFAULT now(), -- UTC
    FOREIGN KEY (user_tg_id) REFERENCES users (tg_id) ON DELETE SET NULL
);

CREATE TABLE messages
(
    id      serial PRIMARY KEY,
    user_tg_id bigint,
    text    text,
    created_at    timestamp with time zone NOT NULL DEFAULT now(),
    FOREIGN KEY (user_tg_id) REFERENCES users (tg_id) ON DELETE SET NULL
);

