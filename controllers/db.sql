CREATE DATABASE db_latihanbackend;

CREATE TABLE accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255),
);

CREATE TABLE games (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    max_player INT,
);  

CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_name VARCHAR(255),
    id_game INT,
    FOREIGN KEY (id_game) REFERENCES games(id),
);

CREATE TABLE participants (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_account INT,
    id_room INT,
    FOREIGN KEY (id_account) REFERENCES accounts(id),
    FOREIGN KEY (id_room) REFERENCES rooms(id),
);

INSERT INTO games (name, max_player) VALUES ('Monster Hunter', 10);
INSERT INTO games (name, max_player) VALUES ('Spider-Man', 5);

INSERT INTO accounts (username) VALUES ('peter');
INSERT INTO accounts (username) VALUES ('john');

INSERT INTO rooms (room_name, id_game) VALUES ('Room 1', 1);
INSERT INTO rooms (room_name, id_game) VALUES ('Room 2', 2);

INSERT INTO participants (id_account, id_room) VALUES (1, 1);
INSERT INTO participants (id_account, id_room) VALUES (2, 2);
```