CREATE TABLE IF NOT EXISTS firebase_tokens
(
    id serial not null unique,
    user_id int references users (id),
    firebase_token text not null
);