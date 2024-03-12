CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   email VARCHAR (300) UNIQUE NOT NULL,
   username VARCHAR(200),
   firstname VARCHAR(100),
   lastname VARCHAR(100)
);
