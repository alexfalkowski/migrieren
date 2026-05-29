CREATE TABLE accounts (
	user_id serial PRIMARY KEY,
	created_at TIMESTAMP NOT NULL
);
INSERT INTO accounts (created_at) VALUES (NOW());
