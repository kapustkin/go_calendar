CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name text
);

CREATE TABLE events (
	id SERIAL PRIMARY KEY, 
	user_id int,
	uuid text,
	start timestamp,
	finish timestamp,
    comment text NULL
);