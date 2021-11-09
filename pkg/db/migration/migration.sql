CREATE TABLE `users`
(
    id   bigint auto_increment,
    user_name varchar(255) NOT NULL,
    user_pass varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
    
);

INSERT INTO `users` (`user_name`,`user_pass`)
VALUES ('user1','123'),
       ('user2','123'),;