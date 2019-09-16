CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name text
);

CREATE TABLE events (
	id SERIAL PRIMARY KEY, 
	eventCreate timestamp,
	user_id int,
	uuid text,
	eventDate timestamp,
    comment text NULL,
	isSended BOOLEAN DEFAULT false
);