package db

var SchemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS users_by_id (
		user_id uuid PRIMARY KEY,
		created_at timestamp,
		email text,
		first_name text,
		last_name text,
		password_hash text,
		username text
	)`,

	`CREATE TABLE IF NOT EXISTS users_by_username (
		username_lower text PRIMARY KEY,
		user_id uuid,
		username text
	)`,

	`CREATE TABLE IF NOT EXISTS friendships (
		user_id uuid,
		friend_id uuid,
		created_at timestamp,
		PRIMARY KEY (user_id, friend_id)
	)`,

	`CREATE TABLE IF NOT EXISTS friend_requests (
		recipient_id uuid,
		sender_id uuid,
		status text,
		created_at timestamp,
		PRIMARY KEY (recipient_id, sender_id)
	)`,

	`CREATE TABLE IF NOT EXISTS channels (
		channel_id uuid PRIMARY KEY,
		name text,
		type text,
		created_by uuid,
		created_at timestamp
	)`,

	`CREATE TABLE IF NOT EXISTS channels_by_user (
		user_id uuid,
		channel_id uuid,
		channel_name text,
		channel_type text,
		joined_at timestamp,
		PRIMARY KEY (user_id, channel_id)
	)`,

	`CREATE TABLE IF NOT EXISTS members_by_channel (
		channel_id uuid,
		user_id uuid,
		role text,
		joined_at timestamp,
		PRIMARY KEY (channel_id, user_id)
	)`,

	`CREATE TABLE IF NOT EXISTS messages_by_channel (
		channel_id uuid,
		bucket int,
		created_at timestamp,
		message_id uuid,
		author_id uuid,
		content text,
		PRIMARY KEY ((channel_id, bucket), created_at, message_id)
	) WITH CLUSTERING ORDER BY (created_at DESC, message_id ASC)`,
}

var TruncateStatements = []string{
	`TRUNCATE users_by_id`,
	`TRUNCATE users_by_username`,
	`TRUNCATE friendships`,
	`TRUNCATE friend_requests`,
	`TRUNCATE channels`,
	`TRUNCATE channels_by_user`,
	`TRUNCATE members_by_channel`,
	`TRUNCATE messages_by_channel`,
}
