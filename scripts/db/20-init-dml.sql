/*
 Start data for the users table
 */
insert into users (username, description, api_key, email, created, disabled)
values ('root', 'root user', '0a2f5aa01a6a61a34607cc9d60a2e996d5d4862ebc3aef53679f559206080e44', 'root@mas2020.me',
        NOW(), false),
       ('andrea', 'simple user', '22599582a225ad6024c572e960371407c9765487e761e9e2c87e1cbbfbde1f27', 'andrea@mas2020.me',
        NOW(), false);