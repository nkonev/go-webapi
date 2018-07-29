CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    name text,
    surname text,
	  lastname text,
	  password text NOT NULL
);
INSERT INTO users VALUES
(DEFAULT, 'root', 's', 'l', 'password'),
(DEFAULT, 'vojtechvitek', '', '', 'password');