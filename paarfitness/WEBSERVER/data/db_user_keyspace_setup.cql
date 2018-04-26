DROP KEYSPACE IF EXISTS user_keyspace;

CREATE KEYSPACE IF NOT EXISTS user_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

USE user_keyspace;

CREATE TABLE IF NOT EXISTS users (
	user_id UUID,
	first_name TEXT,
	last_name TEXT,
	email TEXT,
	pass TEXT,
	country TEXT,
	company TEXT,
	birthday TIMESTAMP,
	created_at TIMESTAMP,
	PRIMARY KEY ((user_id))
);

CREATE TABLE IF NOT EXISTS user_by_email (
	email TEXT,
	pass TEXT,
	user_id UUID,
	PRIMARY KEY ((email))
);

CREATE TABLE IF NOT EXISTS sessions (
	session_id TEXT,
	user_id UUID,
	email TEXT,
	created_at TIMESTAMP,
	device TEXT,
	active BOOLEAN,
	PRIMARY KEY ((session_id))
);

CREATE TABLE IF NOT EXISTS session_by_user_id (
	user_id UUID,
	device TEXT,
	session_id TEXT,
	PRIMARY KEY ((user_id), device)
);

