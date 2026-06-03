CREATE KEYSPACE IF NOT EXISTS my_chat_app
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS my_chat_app.users_by_id (
    user_id uuid PRIMARY KEY,
    created_at timestamp,
    email text,
    first_name text,
    last_name text,
    password_hash text,
    username text
);

CREATE TABLE IF NOT EXISTS my_chat_app.users_by_username (
    username_lower text PRIMARY KEY,
    user_id uuid,
    username text
);
