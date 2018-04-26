DROP KEYSPACE IF EXISTS data_keyspace;

CREATE KEYSPACE IF NOT EXISTS data_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

USE data_keyspace;

CREATE TABLE IF NOT EXISTS objects (
	object_id uuid,
	api text,
	type text,
	data text,
	created_at timestamp,
	PRIMARY KEY ((object_id))
);

CREATE TABLE IF NOT EXISTS results (
	result_id uuid,
	api text,
	type text,
	data text,
	created_at timestamp,
	PRIMARY KEY ((result_id))
);

