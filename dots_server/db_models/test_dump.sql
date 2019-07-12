# blocker
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocker`;

CREATE TABLE `blocker` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `blocker_type` VARCHAR(255) NOT NULL,
  `capacity` int(11) NOT NULL,
  `load` int(11) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_blocker_IDX_LOAD` (`load`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blocker` (`id`, `blocker_type`, `capacity`, `load`, `created`, `updated`)
VALUES
  (1,'Arista-ACL', 100, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (2,'GoBGP-RTBH', 100, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (3,'GoBGP-FlowSpec',  100, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (100,'GoBGP-RTBH',  5, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34');


# blocker_parameters
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
  (1, 1, 'nextHop', '0.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (2, 1, 'host', '127.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (3, 1, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (4, 2, 'nextHop', '0.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (5, 2, 'host', '127.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (6, 2, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (7, 3, 'nextHop', '0.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (8, 3, 'host', '127.0.0.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (9, 3, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (10, 100, 'nextHop', '1.0.0.2', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (11, 100, 'host', '127.1.1.1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (12, 100, 'port', '50051', '2017-04-13 13:44:34', '2017-04-13 13:44:34');


# customer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `customer`;

CREATE TABLE `customer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `common_name` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `customer` (`id`, `common_name`, `created`, `updated`)
VALUES
  (123,'name','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (127,'localhost','2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (128,'client.sample.example.com','2017-04-13 13:44:34','2017-04-13 13:44:34');

# parameter_value
# ------------------------------------------------------------

DROP TABLE IF EXISTS `parameter_value`;

CREATE TABLE `parameter_value` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `type` enum('TARGET_PROTOCOL','FQDN','URI','TRAFFIC_PROTOCOL','ALIAS_NAME') NOT NULL,
  `string_value` varchar(255) DEFAULT NULL,
  `int_value` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `parameter_value` (`id`, `customer_id`, `mitigation_scope_id`, `type`, `string_value`, `int_value`, `created`, `updated`)
VALUES
  (1,123,0,'FQDN','golang.org',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,0,'FQDN','localhost.local',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,128,0,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (4,0,1,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (5,0,2,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34');


# port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `port_range`;

CREATE TABLE `port_range` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `lower_port` int(11) DEFAULT NULL,
  `upper_port` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `port_range` (`id`, `mitigation_scope_id`, `lower_port`, `upper_port`, `created`, `updated`)
VALUES
  (1,1,10000,40000,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,2,10000,65535,'2017-04-13 13:44:34','2017-04-13 13:44:34');

# prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `prefix`;

CREATE TABLE `prefix` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `type` enum('TARGET_PREFIX','SOURCE_IPV4_NETWORK','DESTINATION_IPV4_NETWORK','IP','PREFIX','ADDRESS_RANGE','IP_ADDRESS','TARGET_IP') NOT NULL,
  `addr` varchar(255) DEFAULT NULL,
  `prefix_len` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

INSERT INTO `prefix` (`id`, `customer_id`, `mitigation_scope_id`, `type`, `addr`, `prefix_len`, `created`, `updated`)
VALUES
  (1,123,0,'ADDRESS_RANGE','192.168.1.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,127,0,'ADDRESS_RANGE','129.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (3,127,0,'ADDRESS_RANGE','2003:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (4,127,0,'ADDRESS_RANGE','2003:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (5,128,0,'ADDRESS_RANGE','127.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (6,128,0,'ADDRESS_RANGE','10.100.0.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (7,128,0,'ADDRESS_RANGE','10.101.0.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (8,128,0,'ADDRESS_RANGE','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (9,128,0,'ADDRESS_RANGE','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (10,0,1,'TARGET_IP','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (11,0,1,'TARGET_PREFIX','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (12,0,2,'TARGET_IP','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (13,0,2,'TARGET_PREFIX','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (14,128,0,'ADDRESS_RANGE','1.1.1.69',32,'2017-11-11 20:09:00','2017-11-11 20:09:00'),
  (15,128,0,'ADDRESS_RANGE','1.1.2.0',24,'2017-11-11 20:09:00','2017-11-11 20:09:00');

# prefix trigger when ip address range change
# ------------------------------------------------------------


DROP FUNCTION IF EXISTS MySQLNotification;
CREATE FUNCTION MySQLNotification RETURNS INTEGER SONAME 'mysql-notification.so';

DELIMITER @@

CREATE TRIGGER address_range_trigger AFTER UPDATE ON prefix
FOR EACH ROW
BEGIN
  IF (NEW.type = 'ADDRESS_RANGE') AND (NEW.addr <> OLD.addr OR NEW.prefix_len <> OLD.prefix_len) THEN
    SELECT MySQLNotification('prefix', NEW.customer_id) INTO @x;
  END IF;
END@@

DELIMITER ;


# mitigation_scope
# ------------------------------------------------------------

DROP TABLE IF EXISTS `mitigation_scope`;

CREATE TABLE `mitigation_scope` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `client_identifier` varchar(255) DEFAULT NULL,
  `client_domain_identifier` varchar(255) DEFAULT NULL,
  `mitigation_id` int(11) DEFAULT NULL,
  `status` int(1) DEFAULT NULL,
  `lifetime` int(11) DEFAULT NULL,
  `trigger-mitigation` tinyint(1) DEFAULT NULL,
  `attack-status` int(1) DEFAULT NULL,
  `acl_name` varchar(255) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `mitigation_scope` (`id`, `customer_id`, `client_identifier`, `client_domain_identifier`, `mitigation_id`, `status`, `lifetime`, `trigger-mitigation`,`created`, `updated`)
VALUES
  (1,128,'','',12332,7,1000, 1,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,128,'','',12333,7,1000, 1,'2017-04-13 13:44:34','2017-04-13 13:44:34');

# mitigation_scope trigger when status change
# ------------------------------------------------------------

DELIMITER @@

CREATE TRIGGER status_changed_trigger AFTER UPDATE ON mitigation_scope
FOR EACH ROW
BEGIN
  IF NEW.status <> OLD.status THEN
    SELECT MySQLNotification('mitigation_scope', NEW.id, NEW.customer_id, NEW.client_identifier, NEW.mitigation_id, NEW.status) INTO @x;
  END IF;
END@@

DELIMITER ;


# signal_session_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `signal_session_configuration`;

CREATE TABLE `signal_session_configuration` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) NOT NULL,
  `session_id` int(11) NOT NULL,
  `heartbeat_interval` int(11) DEFAULT NULL,
  `missing_hb_allowed` int(11) DEFAULT NULL,
  `max_retransmit` int(11) DEFAULT NULL,
  `ack_timeout` double DEFAULT NULL,
  `ack_random_factor` double DEFAULT NULL,
  `heartbeat_interval_idle` int(11) DEFAULT NULL,
  `missing_hb_allowed_idle` int(11) DEFAULT NULL,
  `max_retransmit_idle` int(11) DEFAULT NULL,
  `ack_timeout_idle` double DEFAULT NULL,
  `ack_random_factor_idle` double DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_signal_session_configuration_idx_customer_id` (`customer_id`),
  KEY `IDX_signal_session_configuration_idx_session_id` (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# signal_session_configuration trigger when any configuration change
# ------------------------------------------------------------------------------


DELIMITER @@

CREATE TRIGGER session_configuration_changed_trigger AFTER UPDATE ON signal_session_configuration
FOR EACH ROW
BEGIN
  IF (NEW.heartbeat_interval <> OLD.heartbeat_interval) OR (NEW.missing_hb_allowed <> OLD.missing_hb_allowed)
    OR (NEW.max_retransmit <> OLD.max_retransmit) OR (NEW.ack_timeout <> OLD.ack_timeout)
    OR (NEW.ack_random_factor <> OLD.ack_random_factor) OR (NEW.heartbeat_interval_idle <> OLD.heartbeat_interval_idle)
    OR (NEW.missing_hb_allowed_idle <> OLD.missing_hb_allowed_idle) OR (NEW.max_retransmit_idle <> OLD.max_retransmit_idle)
    OR (NEW.ack_timeout_idle <> OLD.ack_timeout_idle) OR (NEW.ack_random_factor_idle <> OLD.ack_random_factor_idle) THEN
    SELECT MySQLNotification('signal_session_configuration', NEW.customer_id, NEW.session_id) INTO @x;
  END IF;
END@@

DELIMITER ;


# protection
# ------------------------------------------------------------

DROP TABLE IF EXISTS `protection`;

CREATE TABLE `protection` (
  `id`                     BIGINT(20)   NOT NULL AUTO_INCREMENT,
  `customer_id`            INT(11)      NOT NULL,
  `target_id`              BIGINT(20)   NOT NULL,
  `target_type`            VARCHAR(255) NOT NULL,
  `acl_name`               VARCHAR(255)          DEFAULT NULL,
  `is_enabled`             TINYINT(1)   NOT NULL,
  `protection_type`        VARCHAR(255) NOT NULL,
  `target_blocker_id`      BIGINT(20)            DEFAULT NULL,
  `started_at`             DATETIME              DEFAULT NULL,
  `finished_at`            DATETIME              DEFAULT NULL,
  `record_time`            DATETIME              DEFAULT NULL,
  `dropped_data_info_id`   BIGINT(20)            DEFAULT NULL,
  `created`                DATETIME              DEFAULT NULL,
  `updated`                DATETIME              DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

insert into `protection` (id, customer_id, target_id, target_type, is_enabled, protection_type, target_blocker_id, started_at, finished_at, record_time, dropped_data_info_id, `created`, `updated`)
VALUES
(100, 128, 1, 'mitigation_request', false, 'RTBH', 1, null, null, null, 1, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(101, 128, 2, 'datachannel_acl', false, 'RTBH', 1, null, null, null, 2, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# gobgp_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `go_bgp_parameter`;

CREATE TABLE `go_bgp_parameter` (
  `id` bigint(20)  NOT NULL AUTO_INCREMENT,
  `protection_id`  BIGINT(20) NOT NULL,
  `target_address` varchar(255) NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into `go_bgp_parameter` (id, protection_id, target_address, `created`, `updated`)
VALUES
(1, 100, '192.168.240.0', '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# protection_status
# ------------------------------------------------------------

DROP TABLE IF EXISTS `protection_status`;

CREATE TABLE `protection_status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bytes_dropped` int(11) DEFAULT NULL,
  `pkts_dropped` int(11) DEFAULT NULL,
  `bps_dropped` int(11) DEFAULT NULL,
  `pps_dropped` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into protection_status (id, bytes_dropped, pkts_dropped, bps_dropped, pps_dropped, created, updated)
VALUES
(1, 0, 0, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(2, 16, 1, 4, 1, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(3, 40, 2, 5, 1, '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(4, 100, 3, 10, 1, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# data_clients
# ------------------------------------------------------------

DROP TABLE IF EXISTS `data_clients`;

CREATE TABLE `data_clients` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `customer_id` INT(11) NOT NULL,
  `cuid` VARCHAR(255) NOT NULL,
  `cdid` VARCHAR(255),
  PRIMARY KEY (`id`),
  KEY `IDX_data_clients_idx_customer_id_cuid` (`customer_id`, `cuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `data_clients` ADD CONSTRAINT UC_dots_clients UNIQUE (`customer_id`, `cuid`);

# data_aliases
# ------------------------------------------------------------

DROP TABLE IF EXISTS `data_aliases`;

CREATE TABLE `data_aliases` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `data_client_id` BIGINT(20) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `valid_through` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_data_aliases_idx_data_client_id_name` (`data_client_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `data_aliases` ADD CONSTRAINT UC_dots_aliases UNIQUE (`data_client_id`, `name`);

# data_acls
# ------------------------------------------------------------

DROP TABLE IF EXISTS `data_acls`;

CREATE TABLE `data_acls` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `data_client_id` BIGINT(20) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `valid_through` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_data_acls_idx_data_client_id_name` (`data_client_id`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `data_acls` ADD CONSTRAINT UC_dots_acls UNIQUE (`data_client_id`, `name`);

# data_acls trigger when activaton_type change
# ------------------------------------------------------------

DELIMITER @@

CREATE TRIGGER activaton_type_changed_trigger AFTER UPDATE ON data_acls
FOR EACH ROW
BEGIN

  DECLARE newContent          VARCHAR(255) DEFAULT NULL;
  DECLARE currentContent      VARCHAR(255) DEFAULT NULL;
  SELECT SUBSTRING_INDEX(NEW.content,",", 3) INTO newContent FROM data_acls limit 1;
  SELECT SUBSTRING_INDEX(OLD.content,",", 3) INTO currentContent FROM data_acls limit 1;

  IF SUBSTRING_INDEX(newContent,"activation-type", -1) <> SUBSTRING_INDEX(currentContent,"activation-type", -1) THEN
    SELECT MySQLNotification('data_acls', NEW.id) INTO @x;
  END IF;
END@@

DELIMITER ;

# arista_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `arista_parameter`;

CREATE TABLE `arista_parameter` (
  `id`                  bigint(20)   NOT NULL AUTO_INCREMENT,
  `protection_id`       bigint(20)   NOT NULL,
  `acl_type`            varchar(255) NOT NULL,
  `acl_filtering_rule`  text     NOT NULL,
  `created`             datetime DEFAULT NULL,
  `updated`             datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# blocker_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocker_configuration`;

CREATE TABLE `blocker_configuration` (
  `id`                bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id`       int(11) NOT NULL,
  `target_type`       VARCHAR(255) NOT NULL,
  `blocker_type`      VARCHAR(255) NOT NULL,
  `created`           datetime DEFAULT NULL,
  `updated`           datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blocker_configuration` (`id`, `customer_id`, `target_type`, `blocker_type`, `created`, `updated`)
VALUES
(1, 128, "mitigation_request", "GoBGP-RTBH", '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
(2, 128, "datachannel_acl", "Arista-ACL", '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# blocker_configuration_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocker_configuration_parameter`;

CREATE TABLE `blocker_configuration_parameter` (
  `id`                       bigint(20) NOT NULL AUTO_INCREMENT,
  `blocker_configuration_id` int(11) NOT NULL,
  `key`                      VARCHAR(255) NOT NULL,
  `value`                    VARCHAR(255) NOT NULL,
  `created`                  datetime DEFAULT NULL,
  `updated`                  datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blocker_configuration_parameter` (`id`, `blocker_configuration_id`, `key`, `value`, `created`, `updated`)
VALUES
  (1, 1, 'vrf', '1.1.1.1:100', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (2, 1, 'aristaConnection', 'arista', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (3, 1, 'aristaInterface', 'Ethernet 1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (4, 2, 'aristaConnection', 'arista', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
  (5, 2, 'aristaInterface', 'Ethernet 1', '2017-04-13 13:44:34', '2017-04-13 13:44:34');


# flow_spec_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `flow_spec_parameter`;

CREATE TABLE `flow_spec_parameter` (
  `id`                  bigint(20)   NOT NULL AUTO_INCREMENT,
  `protection_id`       bigint(20)   NOT NULL,
  `flow_type`           varchar(255) NOT NULL,
  `flow_specification`  text         NOT NULL,
  `created`             datetime DEFAULT NULL,
  `updated`             datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# control_filtering
# ------------------------------------------------------------

DROP TABLE IF EXISTS `control_filtering`;

CREATE TABLE `control_filtering` (
  `id`                  bigint(20)   NOT NULL AUTO_INCREMENT,
  `mitigation_scope_id` bigint(20)   DEFAULT NULL,
  `acl_name`            varchar(255) DEFAULT NULL,
  `created`             datetime     DEFAULT NULL,
  `updated`             datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;