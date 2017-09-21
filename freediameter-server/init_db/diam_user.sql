insert into grp(id, grp_name) values (1, 'group1');

insert into users(id, username, eapmethod, password) values (1, 'client1@example.com', 13, '*');
insert into users(id, username, eapmethod, password) values (2, 'client2@example.com', 13, '*');

insert into user_grp(user, grp) values (1, 1);
insert into user_grp(user, grp) values (2, 1);
