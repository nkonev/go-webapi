CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    name text,
    surname text,
	  lastname text,
	  password text NOT NULL
);
INSERT INTO users VALUES
(0, 'root', 's', 'l', 'password'),
(1, 'vojtechvitek', '', '', 'password');