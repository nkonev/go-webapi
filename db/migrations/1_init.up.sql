CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    email text,
    facebook_id varchar(64),
    creation_type varchar(32) NOT NULL,
	  password text
);
INSERT INTO users(email, password, creation_type) VALUES
('root', 'password', 'email'),
('test@example.com', 'password', 'email');