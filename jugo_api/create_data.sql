DROP TABLE IF EXISTS users;

CREATE TABLE users(
user_id SERIAL PRIMARY KEY,
name VARCHAR(50) UNIQUE NOT NULL,
email VARCHAR(50) UNIQUE NOT NULL,
last_changed TIMESTAMPTZ NOT NULL DEFAULT now(),
file_path TEXT DEFAULT '');

INSERT INTO users(name, email)
VALUES
  ('Blue Train', 'JTrain@mail.com'),
  ('john paul', 'John@gmail.com'),
  ('Gerry cnambodi', 'jerry@gmail.com'),
  ('Sarah Vaughan', 'Sarah@gmail.com');