CREATE TABLE starfleet (
	id INT NOT NULL AUTO_INCREMENT,
	name TEXT NOT NULL,
	passengers INT,
	mission TEXT,
	departure_time_of_ship TIMESTAMP,
	primary key (id)
);

INSERT INTO starfleet (name, passengers, mission, departure_time_of_ship)
VALUES ("USS Enterprise", 1012, NULL, now());
