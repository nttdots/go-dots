# blocker
# ------------------------------------------------------------


INSERT INTO `blocker` (`id`, `type`, `capacity`, `load`, `created`, `updated`)
VALUES
(1,'GoBGP-RTBH', 100, 0,'2017-04-13 13:44:34','2017-04-13 13:44:34');


INSERT INTO `blocker_parameter` (`id`, `blocker_id`, `key`, `value`, `created`, `updated`)
VALUES
(1, 1, 'nextHop', '172.16.236.254','2017-04-13 13:44:34','2017-04-13 13:44:34'),
(2, 1, 'host', '172.16.236.13','2017-04-13 13:44:34','2017-04-13 13:44:34'),
(3, 1, 'port', '50051','2017-04-13 13:44:34','2017-04-13 13:44:34');


# customer
# ------------------------------------------------------------


INSERT INTO `customer` (`id`, `name`, `created`, `updated`)
VALUES
(1,'localhost','2017-04-13 13:44:34','2017-04-13 13:44:34');

# customer_common_name
# ------------------------------------------------------------


INSERT INTO `customer_common_name` (`id`, `customer_id`, `common_name`, `created`, `updated`)
VALUES
(1,1,'client.sample.example.com','2017-04-13 13:44:34','2017-04-13 13:44:34');


# parameter_value
# ------------------------------------------------------------


INSERT INTO `parameter_value` (`id`, `customer_id`, `identifier_id`, `mitigation_scope_id`, `type`, `string_value`, `int_value`, `created`, `updated`)
VALUES
(1,1,0,0,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34');



INSERT INTO `prefix` (`id`, `customer_id`, `identifier_id`, `mitigation_scope_id`, `blocker_id`, `type`, `addr`, `prefix_len`, `created`, `updated`)
VALUES
(1,1,0,0,0,'ADDRESS_RANGE','127.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
(2,1,0,0,0,'ADDRESS_RANGE','2002:db8::',48,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
(3,1,0,0,0,'ADDRESS_RANGE','172.16.236.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34');


