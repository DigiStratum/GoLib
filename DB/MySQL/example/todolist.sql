CREATE DATABASE todolist;
USE todolist;

CREATE TABLE `todo` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`task` VARCHAR(250) DEFAULT NULL,
	`due` DATETIME DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1

INSERT INTO `todo` SET `task` = 'Dishes', `due` = NOW();
INSERT INTO `todo` SET `task` = 'Laundry', `due` = DATE_ADD(NOW(), INTERVAL 1 WEEK);
INSERT INTO `todo` SET `task` = 'Dusting', `due` = DATE_ADD(NOW(), INTERVAL 2 WEEK);
INSERT INTO `todo` SET `task` = 'Sweeping', `due` = DATE_ADD(NOW(), INTERVAL 3 WEEK);

GRANT ALL ON todolist.* TO 'username'@'localhost';
FLUSH PRIVILEGES;