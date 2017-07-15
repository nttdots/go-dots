# テーブルのダンプ blocker
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocker`;

CREATE TABLE `blocker` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(20) NOT NULL,
  `capacity` int(11) NOT NULL,
  `load` int(11) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_blocker_IDX_LOAD` (`load`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blocker` (`id`, `type`, `capacity`, `load`, `created`, `updated`)
VALUES
  (1,'GoBGP-RTBH', 100, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (2,'GoBGP-RTBH', 100, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (3,'GoBGP-RTBH',  10, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (100, 'GoBGP-RTBH',  5, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34');


# テーブルのダンプ blocker_parameters
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocker_parameter`;

CREATE TABLE `blocker_parameter` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `blocker_id` bigint(20) NOT NULL,
  `key` VARCHAR(255) NOT NULL,
  `value` VARCHAR(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blocker_parameter` (`id`, `blocker_id`, `key`, `value`, `created`, `updated`)
VALUES
  (1, 1, 'nextHop', '0.0.0.0', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (2, 1, 'host', '127.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (3, 1, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (4, 2, 'nextHop', '0.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (5, 2, 'host', '127.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (6, 2, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (7, 3, 'nextHop', '0.0.0.2', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (8, 3, 'host', '127.0.2.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (9, 3, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (10, 100, 'nextHop', '1.0.0.2', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (11, 100, 'host', '127.1.1.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (12, 100, 'port', '50056', '2017-04-13 13:44:34', '2017-04-13 13:44:34');


# テーブルのダンプ customer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `customer`;

CREATE TABLE `customer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `customer` (`id`, `name`, `created`, `updated`)
VALUES
  (123,'name','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (127,'localhost','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (128,'client.sample.example.com','2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ customer_common_name
# ------------------------------------------------------------

DROP TABLE IF EXISTS `customer_common_name`;

CREATE TABLE `customer_common_name` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `common_name` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_customer_common_name_IDX_CUSTOMER_ID` (`customer_id`),
  KEY `IDX_customer_common_name_IDX_COMMON_NAME` (`common_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `customer_common_name` (`id`, `customer_id`, `common_name`, `created`, `updated`)
VALUES
  (1,123,'commonName','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,'local-host','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,128,'client.sample.example.com','2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ identifier
# ------------------------------------------------------------

DROP TABLE IF EXISTS `identifier`;

CREATE TABLE `identifier` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `alias_name` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_identifier_IDX_CUSTOMER_ID` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# テーブルのダンプ login_profile
# ------------------------------------------------------------

DROP TABLE IF EXISTS `login_profile`;

CREATE TABLE `login_profile` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `blocker_id` bigint(20) NOT NULL,
  `login_method` varchar(255) NOT NULL,
  `login_name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_login_profile_IDX_BLOCKER_ID` (`blocker_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `login_profile` (`id`, `blocker_id`, `login_method`, `login_name`, `password`, `created`, `updated`)
VALUES
  (1,123,'ssh','go','receiver192.168.10.20','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,'ssh','go','receiver192.168.10.30','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,128,'ssh','go','receiver192.168.10.40','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (4,100,'ssh','go','receiver192.168.100.40','2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ parameter_value
# ------------------------------------------------------------

DROP TABLE IF EXISTS `parameter_value`;

CREATE TABLE `parameter_value` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `identifier_id` bigint(20) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `type` enum('TARGET_PROTOCOL','FQDN','URI','E_164','TRAFFIC_PROTOCOL','ALIAS') NOT NULL,
  `string_value` varchar(255) DEFAULT NULL,
  `int_value` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `parameter_value` (`id`, `customer_id`, `identifier_id`, `mitigation_scope_id`, `type`, `string_value`, `int_value`, `created`, `updated`)
VALUES
  (1,123,0,0,'FQDN','golang.org',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,0,0,'FQDN','localhost.local',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,128,0,0,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (4,0,0,1,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (5,0,0,2,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `port_range`;

CREATE TABLE `port_range` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `identifier_id` bigint(20) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `lower_port` int(11) DEFAULT NULL,
  `upper_port` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `port_range` (`id`, `identifier_id`, `mitigation_scope_id`, `lower_port`, `upper_port`, `created`, `updated`)
VALUES
  (1,0,1,10000,40000,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,0,2,10000,65535,'2017-04-13 13:44:34','2017-04-13 13:44:34');

# テーブルのダンプ prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `prefix`;

CREATE TABLE `prefix` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `identifier_id` bigint(20) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `blocker_id` bigint(20) DEFAULT NULL,
  `access_control_list_entry_id` bigint(20) DEFAULT NULL,
  `type` enum('TARGET_PREFIX','SOURCE_IPV4_NETWORK','DESTINATION_IPV4_NETWORK','IP','PREFIX','ADDRESS_RANGE','IP_ADDRESS','TARGET_IP') NOT NULL,
  `addr` varchar(255) DEFAULT NULL,
  `prefix_len` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

INSERT INTO `prefix` (`id`, `customer_id`, `identifier_id`, `mitigation_scope_id`, `blocker_id`, `type`, `addr`, `prefix_len`, `created`, `updated`)
VALUES
  (1,123,0,0,0,'ADDRESS_RANGE','192.168.1.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,0,0,0,'ADDRESS_RANGE','129.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,127,0,0,0,'ADDRESS_RANGE','2003:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (4,127,0,0,0,'ADDRESS_RANGE','2003:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (5,128,0,0,0,'ADDRESS_RANGE','127.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (6,128,0,0,0,'ADDRESS_RANGE','10.100.0.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (7,128,0,0,0,'ADDRESS_RANGE','10.101.0.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (8,128,0,0,0,'ADDRESS_RANGE','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (9,128,0,0,0,'ADDRESS_RANGE','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (10,0,0,1,0,'TARGET_IP','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (11,0,0,1,0,'TARGET_PREFIX','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (12,0,0,2,0,'TARGET_IP','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (13,0,0,2,0,'TARGET_PREFIX','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ mitigation_scope
# ------------------------------------------------------------

DROP TABLE IF EXISTS `mitigation_scope`;

CREATE TABLE `mitigation_scope` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `mitigation_id` int(11) DEFAULT NULL,
  `lifetime` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `mitigation_scope` (`id`, `customer_id`, `mitigation_id`, `lifetime`, `created`, `updated`)
VALUES
  (1,128,12332,1000,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,128,12333,1000,'2017-04-13 13:44:34','2017-04-13 13:44:34');


# テーブルのダンプ signal_session_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `signal_session_configuration`;

CREATE TABLE `signal_session_configuration` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `session_id` int(11) NOT NULL,
  `heartbeat_interval` int(11) DEFAULT NULL,
  `max_retransmit` int(11) DEFAULT NULL,
  `ack_timeout` int(11) DEFAULT NULL,
  `ack_random_factor` double DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_signal_session_configuration_idx_customer_id` (`customer_id`),
  KEY `IDX_signal_session_configuration_idx_session_id` (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# テーブルのダンプ protection
# ------------------------------------------------------------

DROP TABLE IF EXISTS `protection`;

CREATE TABLE `protection` (
  `id`                     BIGINT(20)   NOT NULL AUTO_INCREMENT,
  `mitigation_id`          INT(11)      NOT NULL,
  `is_enabled`             TINYINT(1)   NOT NULL,
  `type`                   VARCHAR(255) NOT NULL,
  `target_blocker_id`      BIGINT(20)            DEFAULT NULL,
  `started_at`             DATETIME              DEFAULT NULL,
  `finished_at`            DATETIME              DEFAULT NULL,
  `record_time`            DATETIME              DEFAULT NULL,
  `forwarded_data_info_id` BIGINT(20)            DEFAULT NULL,
  `blocked_data_info_id`   BIGINT(20)            DEFAULT NULL,
  `created`                DATETIME              DEFAULT NULL,
  `updated`                DATETIME              DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_protection_idx_mitigation_id` (mitigation_id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

insert into `protection` (id, mitigation_id, is_enabled, `type`, target_blocker_id, started_at, finished_at, record_time, forwarded_data_info_id, blocked_data_info_id, `created`, `updated`)
VALUES
(100, 1, false, 'RTBH', 1, null, null, null, 1, 2, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(101, 2, false, 'RTBH', 1, null, null, null, 3, 4, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# テーブルのダンプ protection_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `protection_parameter`;

CREATE TABLE `protection_parameter` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `protection_id` BIGINT(20) NOT NULL,
  `key` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into `protection_parameter` (id, protection_id, `key`, `value`, `created`, `updated`)
VALUES
(1, 100, 'customerId', '1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(2, 100, 'target', '192.168.240.0', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(3, 101, 'customerId', '10', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(4, 101, 'target', '192.168.241.0', '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# テーブルのダンプ protection_status
# ------------------------------------------------------------

DROP TABLE IF EXISTS `protection_status`;

CREATE TABLE `protection_status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `total_packets` int(11) DEFAULT NULL,
  `total_bits` int(11) DEFAULT NULL,
  `peak_throughput_id` bigint(20) DEFAULT NULL,
  `average_throughput_id` bigint(20) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into protection_status (id, total_packets, total_bits, peak_throughput_id, average_throughput_id, created, updated)
VALUES
(1, 0, 0, 1, 2, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(2, 0, 0, 3, 4, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(3, 0, 0, 5, 6, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(4, 0, 0, 7, 8, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# テーブルのダンプ throughput_data
# ------------------------------------------------------------

DROP TABLE IF EXISTS `throughput_data`;

CREATE TABLE `throughput_data` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pps` int(11) DEFAULT NULL,
  `bps` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into throughput_data (id, pps, bps, created, updated)
values
(1, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(2, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(3, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(4, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(5, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(6, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(7, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(8, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# テーブルのダンプ access_control_list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `access_control_list`;

CREATE TABLE `access_control_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_access_control_list_idx_customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into access_control_list(id, customer_id, name, type, created, updated)
values
(1, 127, 'abc', 'abc', '2017-06-13 12:00:00', '2017-06-14 15:00:00');

# テーブルのダンプ access_control_list_entry
# ------------------------------------------------------------

DROP TABLE IF EXISTS `access_control_list_entry`;

CREATE TABLE `access_control_list_entry` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `access_control_list_id` bigint(20) NOT NULL,
  `rule_name` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_access_control_list_entry_idx_access_control_list_id` (`access_control_list_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into access_control_list_entry(id, access_control_list_id, rule_name, created, updated)
VALUES
(1, 1, 'abc', '2017-06-13 12:00:00', '2017-06-14 15:00:00');

# テーブルのダンプ acl_rule_action
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acl_rule_action`;

CREATE TABLE `acl_rule_action` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `access_control_list_entry_id` bigint(20) NOT NULL,
  `type` enum('DENY','PERMIT','RATE_LIMIT') NOT NULL,
  `action` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_acl_rule_action_idx_access_control_list_entry_id` (`access_control_list_entry_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
