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

-- Incoming requests (partition by recipient so "my pending" is a single query)
CREATE TABLE friend_requests (
    recipient_id uuid,
    sender_id uuid,
    status text,        -- 'pending' | 'accepted' | 'rejected'
    created_at timestamp,
    PRIMARY KEY (recipient_id, sender_id)
);

-- Mutual friendships, stored bidirectionally
CREATE TABLE friendships (
    user_id uuid,
    friend_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, friend_id)
);

-- Channel metadata lookup
CREATE TABLE channels (
    channel_id  uuid,
    name        text,
    type        text,            -- 'public' | 'group' | 'dm'
    created_by  uuid,
    created_at  timestamp,
    PRIMARY KEY (channel_id)
);

-- "What channels am I in?" (sidebar) — partitioned by user
CREATE TABLE channels_by_user (
    user_id       uuid,
    channel_id    uuid,
    channel_name  text,          -- denormalized so the sidebar is ONE read
    channel_type  text,
    joined_at     timestamp,
    PRIMARY KEY (user_id, channel_id)
);

-- "Who's in this channel?" (member list, fan-out) — partitioned by channel
CREATE TABLE members_by_channel (
    channel_id  uuid,
    user_id     uuid,
    role        text,            -- 'owner' | 'admin' | 'member'
    joined_at   timestamp,
    PRIMARY KEY (channel_id, user_id)
);

-- Recent messages in a channel
CREATE TABLE messages_by_channel (
    channel_id  uuid,
    bucket      int,             -- e.g. 20260618 (UTC day)
    message_id  timeuuid,        -- carries time + ordering + uniqueness
    author_id   uuid,
    content     text,
    PRIMARY KEY ((channel_id, bucket), message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);