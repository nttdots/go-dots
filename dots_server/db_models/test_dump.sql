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
  `type` enum('TARGET_PORT','SOURCE_PORT') NOT NULL,
  `lower_port` int(11) DEFAULT NULL,
  `upper_port` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `port_range` (`id`, `mitigation_scope_id`, `type`, `lower_port`, `upper_port`, `created`, `updated`)
VALUES
  (1,1,'TARGET_PORT',10000,40000,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,2,'SOURCE_PORT',10000,65535,'2017-04-13 13:44:34','2017-04-13 13:44:34');

# icmp_type_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `icmp_type_range`;

CREATE TABLE `icmp_type_range` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `lower_type` int(11) DEFAULT NULL,
  `upper_type` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `icmp_type_range` (`id`, `mitigation_scope_id`, `lower_type`, `upper_type`, `created`, `updated`)
VALUES
  (1,1,10,11,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
  (2,2,12,13,'2017-04-13 13:44:34','2017-04-13 13:44:34');

# prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `prefix`;

CREATE TABLE `prefix` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `customer_id` int(11) DEFAULT NULL,
  `mitigation_scope_id` bigint(20) DEFAULT NULL,
  `type` enum('TARGET_PREFIX','SOURCE_PREFIX','SOURCE_IPV4_NETWORK','DESTINATION_IPV4_NETWORK','IP','PREFIX','ADDRESS_RANGE','IP_ADDRESS','TARGET_IP') NOT NULL,
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
    SELECT MySQLNotification('mitigation_scope', NEW.id, NEW.customer_id, NEW.client_identifier, NEW.mitigation_id, NEW.client_domain_identifier, NEW.status) INTO @x;
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

# telemetry_setup
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_setup`;

CREATE TABLE `telemetry_setup` (
  `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
  `customer_id` int(11)      NOT NULL,
  `cuid`        varchar(255) NOT NULL,
  `cdid`        varchar(255) DEFAULT NULL,
  `tsid`        int(11)      NOT NULL,
  `setup_type`  enum('TELEMETRY_CONFIGURATION','PIPE','BASELINE') NOT NULL,
  `created`     datetime     DEFAULT NULL,
  `updated`     datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_configuration`;

CREATE TABLE `telemetry_configuration` (
  `id`                           bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_setup_id`                bigint(20)   NOT NULL,
  `measurement_interval`         enum('HOUR','DAY','WEEK','MONTH') NOT NULL,
  `measurement_sample`           enum('SECOND','5_SECONDS','30_SECONDS','ONE_MINUTE','5_MINUTES','10_MINUTES','30_MINUTES','ONE_HOUR') NOT NULL,
  `low_percentile`               double       DEFAULT NULL,
  `mid_percentile`               double       DEFAULT NULL,
  `high_percentile`              double       DEFAULT NULL,
  `server_originated_telemetry`  tinyint(1)   NOT NULL,
  `telemetry_notify_interval`    int(11)      DEFAULT NULL,
  `created`                      datetime     DEFAULT NULL,
  `updated`                      datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# unit_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `unit_configuration`;

CREATE TABLE `unit_configuration` (
  `id`             bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_config_id` bigint(20) NOT NULL,
  `unit`           enum('PACKETS_PS','BITS_PS','BYTES_PS') NOT NULL,
  `unit_status`    tinyint(1) DEFAULT NULL,
  `created`        datetime   DEFAULT NULL,
  `updated`        datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# total_pipe_capability
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_pipe_capacity`;

CREATE TABLE `total_pipe_capacity` (
  `id`            bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_setup_id` bigint(20)   NOT NULL,
  `link_id`       varchar(255) DEFAULT NULL,
  `capacity`      int(11)      DEFAULT NULL,
  `unit`          enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `created`       datetime     DEFAULT NULL,
  `updated`       datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# baseline
# ------------------------------------------------------------

DROP TABLE IF EXISTS `baseline`;

CREATE TABLE `baseline` (
  `id`            bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_setup_id` bigint(20)   NOT NULL,
  `baseline_id`   int(11)      NOT NULL,
  `created`       datetime     DEFAULT NULL,
  `updated`       datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_prefix`;

CREATE TABLE `telemetry_prefix` (
  `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`        enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`     bigint(20)   NOT NULL,
  `prefix_type` enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `addr`        varchar(255) DEFAULT NULL,
  `prefix_len`  int(11)      DEFAULT NULL,
  `created`     datetime     DEFAULT NULL,
  `updated`     datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

# telemetry_port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_port_range`;

CREATE TABLE `telemetry_port_range` (
  `id`          bigint(20) NOT NULL AUTO_INCREMENT,
  `type`        enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`     bigint(20) NOT NULL,
  `prefix_type` enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `lower_port`  int(11)    NOT NULL,
  `upper_port`  int(11)    DEFAULT NULL,
  `created`     datetime   DEFAULT NULL,
  `updated`     datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_parameter_value
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_parameter_value`;

CREATE TABLE `telemetry_parameter_value` (
  `id`             bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`           enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`        bigint(20)   NOT NULL,
  `parameter_type` enum('TARGET_PROTOCOL','FQDN','URI','ALIAS_NAME') NOT NULL,
  `string_value`   varchar(255) DEFAULT NULL,
  `int_value`      int(11)      DEFAULT NULL,
  `created`        datetime     DEFAULT NULL,
  `updated`        datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# traffic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `traffic`;

CREATE TABLE `traffic` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`              enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `type_id`           bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# traffic_per_protocol
# ------------------------------------------------------------

DROP TABLE IF EXISTS `traffic_per_protocol`;

CREATE TABLE `traffic_per_protocol` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`              enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`           bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `protocol`          int(11)     NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# traffic_per_port
# ------------------------------------------------------------

DROP TABLE IF EXISTS `traffic_per_port`;

CREATE TABLE `traffic_per_port` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`              enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`           bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `port`              int(11)     NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# total_connection_capacity
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_connection_capacity`;

CREATE TABLE `total_connection_capacity` (
  `id`                        bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_baseline_id`          bigint(20)   NOT NULL,
  `protocol`                  int(11)      NOT NULL,
  `connection`                int(11)      DEFAULT NULL,
  `connection_client`         int(11)      DEFAULT NULL,
  `embryonic`                 int(11)      DEFAULT NULL,
  `embryonic_client`          int(11)      DEFAULT NULL,
  `connection_ps`             int(11)      DEFAULT NULL,
  `connection_client_ps`      int(11)      DEFAULT NULL,
  `request_ps`                int(11)      DEFAULT NULL,
  `request_client_ps`         int(11)      DEFAULT NULL,
  `partial_request_ps`        int(11)      DEFAULT NULL,
  `partial_request_client_ps` int(11)      DEFAULT NULL,
  `created`                   datetime     DEFAULT NULL,
  `updated`                   datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# total_connection_capacity_per_port
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_connection_capacity_per_port`;

CREATE TABLE `total_connection_capacity_per_port` (
  `id`                        bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_baseline_id`          bigint(20)   NOT NULL,
  `protocol`                  int(11)      NOT NULL,
  `port`                      int(11)      NOT NULL,
  `connection`                int(11)      DEFAULT NULL,
  `connection_client`         int(11)      DEFAULT NULL,
  `embryonic`                 int(11)      DEFAULT NULL,
  `embryonic_client`          int(11)      DEFAULT NULL,
  `connection_ps`             int(11)      DEFAULT NULL,
  `connection_client_ps`      int(11)      DEFAULT NULL,
  `request_ps`                int(11)      DEFAULT NULL,
  `request_client_ps`         int(11)      DEFAULT NULL,
  `partial_request_ps`        int(11)      DEFAULT NULL,
  `partial_request_client_ps` int(11)      DEFAULT NULL,
  `created`                   datetime     DEFAULT NULL,
  `updated`                   datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_pre_mitigation
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_pre_mitigation`;

CREATE TABLE `telemetry_pre_mitigation` (
  `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
  `customer_id` int(11)      NOT NULL,
  `cuid`        varchar(255) NOT NULL,
  `cdid`        varchar(255) DEFAULT NULL,
  `tmid`        int(11)      NOT NULL,
  `created`     datetime     DEFAULT NULL,
  `updated`     datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_pre_mitigation trigger when any attribute of telemetry_pre_mitigation change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER telemetry_pre_mitigation_trigger AFTER UPDATE ON telemetry_pre_mitigation
FOR EACH ROW
BEGIN
  IF NEW.updated <> OLD.updated THEN
    SELECT MySQLNotification('telemetry_pre_mitigation', NEW.id) INTO @x;
  END IF;
END@@
DELIMITER ;

# total_attack_connection
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_attack_connection`;

CREATE TABLE `total_attack_connection` (
  `id`                bigint(20) NOT NULL AUTO_INCREMENT,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `prefix_type_id`    bigint(20) NOT NULL,
  `percentile_type`   enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') NOT NULL,
  `protocol`          int(11)  NOT NULL,
  `connection`        int(11)  DEFAULT NULL,
  `embryonic`         int(11)  DEFAULT NULL,
  `connection_ps`     int(11)  DEFAULT NULL,
  `request_ps`        int(11)  DEFAULT NULL,
  `partial_request_ps`int(11)  DEFAULT NULL,
  `created`           datetime DEFAULT NULL,
  `updated`           datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# total_attack_connection_port
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_attack_connection_port`;

CREATE TABLE `total_attack_connection_port` (
  `id`                     bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20) NOT NULL,
  `percentile_type`        enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') NOT NULL,
  `protocol`               int(11)  NOT NULL,
  `port`                   int(11)  NOT NULL,
  `connection`             int(11)  DEFAULT NULL,
  `embryonic`              int(11)  DEFAULT NULL,
  `connection_ps`          int(11)  DEFAULT NULL,
  `request_ps`             int(11)  DEFAULT NULL,
  `partial_request_ps`     int(11)  DEFAULT NULL,
  `created`                datetime DEFAULT NULL,
  `updated`                datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# attack_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `attack_detail`;

CREATE TABLE `attack_detail` (
  `id`                     bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20),
  `vendor_id`              int(11)      NOT NULL,
  `attack_id`              int(11) NOT NULL,
  `attack_name`            varchar(255),
  `attack_severity`        enum('NONE','LOW','MEDIUM','HIGH','UNKNOWN') NOT NULL,
  `start_time`             int(11),
  `end_time`               int(11),
  `created`                datetime DEFAULT NULL,
  `updated`                datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# source_count
# ------------------------------------------------------------

DROP TABLE IF EXISTS `source_count`;

CREATE TABLE `source_count` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `low_percentile_g`      int(11),
  `mid_percentile_g`      int(11),
  `high_percentile_g`     int(11),
  `peak_g`                int(11),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# top_talker
# ------------------------------------------------------------

DROP TABLE IF EXISTS `top_talker`;

CREATE TABLE `top_talker` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `spoofed_status`        tinyint(1),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_icmp_type_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_icmp_type_range`;

CREATE TABLE `telemetry_icmp_type_range` (
  `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20) NOT NULL,
  `lower_type`         int(11)    NOT NULL,
  `upper_type`         int(11)    DEFAULT NULL,
  `created`            datetime   DEFAULT NULL,
  `updated`            datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_traffic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_traffic`;

CREATE TABLE `telemetry_traffic` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `prefix_type_id`    bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `protocol`          int(11)     NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `telemetry_traffic` (`id`, `prefix_type`, `prefix_type_id`, `traffic_type`, `unit`, `protocol`, `low_percentile_g`, `mid_percentile_g`, `high_percentile_g`, `peak_g`, `created`, `updated`)
VALUES
  (1, 'TARGET_PREFIX', 1, 'TOTAL_TRAFFIC', 'PACKETS_PS', 6, 0, 100, 0, 0, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# telemetry_total_attack_connection
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_total_attack_connection`;

CREATE TABLE `telemetry_total_attack_connection` (
  `id`                bigint(20) NOT NULL AUTO_INCREMENT,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `prefix_type_id`    bigint(20) NOT NULL,
  `percentile_type`   enum('LOW_PERCENTILE_C','MID_PERCENTILE_C','HIGH_PERCENTILE_C','PEAK_C') NOT NULL,
  `connection`        int(11)  DEFAULT NULL,
  `embryonic`         int(11)  DEFAULT NULL,
  `connection_ps`     int(11)  DEFAULT NULL,
  `request_ps`        int(11)  DEFAULT NULL,
  `partial_request_ps`int(11)  DEFAULT NULL,
  `created`           datetime DEFAULT NULL,
  `updated`           datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `telemetry_total_attack_connection` (`id`, `prefix_type`, `prefix_type_id`, `percentile_type`, `connection`, `embryonic`, `connection_ps`, `request_ps`, `partial_request_ps`, `created`, `updated`)
VALUES
  (1, 'TARGET_PREFIX', 1, 'LOW_PERCENTILE_C', 200, 201, 202, 203, 204, '2017-04-13 13:44:34', '2017-04-13 13:44:34');

# telemetry_attack_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_attack_detail`;

CREATE TABLE `telemetry_attack_detail` (
  `id`                  bigint(20) NOT NULL AUTO_INCREMENT,
  `mitigation_scope_id` bigint(20) NOT NULL,
  `vendor_id`           int(11)    NOT NULL,
  `attack_id`           int(11)    NOT NULL,
  `attack_name`         varchar(255),
  `attack_severity`     enum('NONE','LOW','MEDIUM','HIGH','UNKNOWN') NOT NULL,
  `start_time`          int(11),
  `end_time`            int(11),
  `created`             datetime DEFAULT NULL,
  `updated`             datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_attack_detail trigger when any attribute of telemetry_attack_detail change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER telemetry_attack_detail_trigger AFTER UPDATE ON telemetry_attack_detail
FOR EACH ROW
BEGIN
  IF NEW.updated <> OLD.updated THEN
    SELECT MySQLNotification('telemetry_attack_detail', NEW.mitigation_scope_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# telemetry_source_count
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_source_count`;

CREATE TABLE `telemetry_source_count` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `low_percentile_g`      int(11),
  `mid_percentile_g`      int(11),
  `high_percentile_g`     int(11),
  `peak_g`                int(11),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_top_talker
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_top_talker`;

CREATE TABLE `telemetry_top_talker` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `spoofed_status`        tinyint(1),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_source_prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_source_prefix`;

CREATE TABLE `telemetry_source_prefix` (
  `id`                 bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20)   NOT NULL,
  `addr`               varchar(255) DEFAULT NULL,
  `prefix_len`         int(11)      DEFAULT NULL,
  `created`            datetime     DEFAULT NULL,
  `updated`            datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;


# telemetry_source_port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_source_port_range`;

CREATE TABLE `telemetry_source_port_range` (
  `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20) NOT NULL,
  `lower_port`         int(11)    NOT NULL,
  `upper_port`         int(11)    DEFAULT NULL,
  `created`            datetime   DEFAULT NULL,
  `updated`            datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# telemetry_source_icmp_type_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_source_icmp_type_range`;

CREATE TABLE `telemetry_source_icmp_type_range` (
  `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20) NOT NULL,
  `lower_type`         int(11)    NOT NULL,
  `upper_type`         int(11)    DEFAULT NULL,
  `created`            datetime   DEFAULT NULL,
  `updated`            datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_telemetry_pre_mitigation
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_telemetry_pre_mitigation`;

CREATE TABLE `uri_filtering_telemetry_pre_mitigation` (
  `id`              bigint(20)   NOT NULL AUTO_INCREMENT,
  `customer_id`     int(11)      NOT NULL,
  `cuid`            varchar(255) NOT NULL,
  `cdid`            varchar(255) DEFAULT NULL,
  `tmid`            int(11)      NOT NULL,
  `target_prefix`   varchar(255) NOT NULL,
  `lower_port`      int(11)      NOT NULL,
  `upper_port`      int(11)      NOT NULL,
  `target_protocol` int(11)      NOT NULL,
  `target_fqdn`     varchar(255) NOT NULL,
  `alias_name`      varchar(255) NOT NULL,
  `created`         datetime     DEFAULT NULL,
  `updated`         datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_traffic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_traffic`;

CREATE TABLE `uri_filtering_traffic` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `prefix_type_id`    bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_traffic trigger when any attribute of uri_filtering_traffic change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_traffic_trigger AFTER UPDATE ON uri_filtering_traffic
FOR EACH ROW
BEGIN
  IF NEW.unit <> OLD.unit OR NEW.low_percentile_g <> OLD.low_percentile_g OR NEW.mid_percentile_g <> OLD.mid_percentile_g
     OR NEW.high_percentile_g <> OLD.high_percentile_g OR NEW.peak_g <> OLD.peak_g THEN
    SELECT MySQLNotification('uri_filtering_traffic', NEW.prefix_type, NEW.prefix_type_id) INTO @x;
  END IF;
END@@
DELIMITER ;


# uri_filtering_traffic_per_protocol
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_traffic_per_protocol`;

CREATE TABLE `uri_filtering_traffic_per_protocol` (
  `id`                     bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20)   NOT NULL,
  `traffic_type`           enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`                   enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `protocol`               int(11)     NOT NULL,
  `low_percentile_g`       int(11)     DEFAULT NULL,
  `mid_percentile_g`       int(11)     DEFAULT NULL,
  `high_percentile_g`      int(11)     DEFAULT NULL,
  `peak_g`                 int(11)     DEFAULT NULL,
  `created`                datetime    DEFAULT NULL,
  `updated`                datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_traffic_per_protocol trigger when any attribute of uri_filtering_traffic_per_protocol change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_traffic_per_protocol_trigger AFTER UPDATE ON uri_filtering_traffic_per_protocol
FOR EACH ROW
BEGIN
  IF NEW.unit <> OLD.unit OR NEW.protocol <> OLD.protocol OR NEW.low_percentile_g <> OLD.low_percentile_g OR NEW.mid_percentile_g <> OLD.mid_percentile_g 
    OR NEW.high_percentile_g <> OLD.high_percentile_g OR NEW.peak_g <> OLD.peak_g THEN
    SELECT MySQLNotification('uri_filtering_traffic_per_protocol', NEW.tele_pre_mitigation_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_traffic_per_port
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_traffic_per_port`;

CREATE TABLE `uri_filtering_traffic_per_port` (
  `id`                     bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20)   NOT NULL,
  `traffic_type`           enum('TOTAL_TRAFFIC_NORMAL','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`                   enum('PACKETS_PS','BITS_PS','BYTES_PS','KILOPACKETS_PS','KILOBITS_PS','KILOBYTES_PS','MEGAPACKETS_PS','MEGABITS_PS','MEGABYTES_PS','GIGAPACKETS_PS','GIGABITS_PS','GIGABYTES_PS','TERAPACKETS_PS','TERABITS_PS','TERABYTES_PS') NOT NULL,
  `port`                   int(11)     NOT NULL,
  `low_percentile_g`       int(11)     DEFAULT NULL,
  `mid_percentile_g`       int(11)     DEFAULT NULL,
  `high_percentile_g`      int(11)     DEFAULT NULL,
  `peak_g`                 int(11)     DEFAULT NULL,
  `created`                datetime    DEFAULT NULL,
  `updated`                datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_traffic_per_port trigger when any attribute of uri_filtering_traffic_per_port change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_traffic_per_port_trigger AFTER UPDATE ON uri_filtering_traffic_per_port
FOR EACH ROW
BEGIN
  IF NEW.unit <> OLD.unit OR NEW.port <> OLD.port OR NEW.low_percentile_g <> OLD.low_percentile_g OR NEW.mid_percentile_g <> OLD.mid_percentile_g 
    OR NEW.high_percentile_g <> OLD.high_percentile_g OR NEW.peak_g <> OLD.peak_g THEN
    SELECT MySQLNotification('uri_filtering_traffic_per_port', NEW.tele_pre_mitigation_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_total_attack_connection
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_total_attack_connection`;

CREATE TABLE `uri_filtering_total_attack_connection` (
  `id`                bigint(20) NOT NULL AUTO_INCREMENT,
  `prefix_type`       enum('TARGET_PREFIX','SOURCE_PREFIX') NOT NULL,
  `prefix_type_id`    bigint(20) NOT NULL,
  `percentile_type`   enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') NOT NULL,
  `protocol`          int(11)  NOT NULL,
  `connection`        int(11)  DEFAULT NULL,
  `embryonic`         int(11)  DEFAULT NULL,
  `connection_ps`     int(11)  DEFAULT NULL,
  `request_ps`        int(11)  DEFAULT NULL,
  `partial_request_ps`int(11)  DEFAULT NULL,
  `created`           datetime DEFAULT NULL,
  `updated`           datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_total_attack_connection trigger when any attribute of uri_filtering_total_attack_connection change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_total_attack_connection_trigger AFTER UPDATE ON uri_filtering_total_attack_connection
FOR EACH ROW
BEGIN
  IF NEW.protocol <> OLD.protocol OR NEW.connection <> OLD.connection OR NEW.embryonic <> OLD.embryonic OR NEW.connection_ps <> OLD.connection_ps
     OR NEW.request_ps <> OLD.request_ps OR NEW.partial_request_ps <> OLD.partial_request_ps THEN
    SELECT MySQLNotification('uri_filtering_total_attack_connection', NEW.prefix_type, NEW.prefix_type_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_total_attack_connection_port
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_total_attack_connection_port`;

CREATE TABLE `uri_filtering_total_attack_connection_port` (
  `id`                     bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20) NOT NULL,
  `percentile_type`        enum('LOW_PERCENTILE_L','MID_PERCENTILE_L','HIGH_PERCENTILE_L','PEAK_L') NOT NULL,
  `protocol`               int(11)  NOT NULL,
  `port`                   int(11)  NOT NULL,
  `connection`             int(11)  DEFAULT NULL,
  `embryonic`              int(11)  DEFAULT NULL,
  `connection_ps`          int(11)  DEFAULT NULL,
  `request_ps`             int(11)  DEFAULT NULL,
  `partial_request_ps`     int(11)  DEFAULT NULL,
  `created`                datetime DEFAULT NULL,
  `updated`                datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_total_attack_connection_port trigger when any attribute of uri_filtering_total_attack_connection_port change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_total_attack_connection_port_trigger AFTER UPDATE ON uri_filtering_total_attack_connection_port
FOR EACH ROW
BEGIN
  IF NEW.protocol <> OLD.protocol OR NEW.port <> OLD.port OR NEW.connection <> OLD.connection OR NEW.embryonic <> OLD.embryonic OR NEW.connection_ps <> OLD.connection_ps
     OR NEW.request_ps <> OLD.request_ps OR NEW.partial_request_ps <> OLD.partial_request_ps THEN
    SELECT MySQLNotification('uri_filtering_total_attack_connection_port', NEW.tele_pre_mitigation_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_attack_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_attack_detail`;

CREATE TABLE `uri_filtering_attack_detail` (
  `id`                     bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_pre_mitigation_id` bigint(20),
  `vendor_id`              int(11)      NOT NULL,
  `attack_id`              int(11)      NOT NULL,
  `attack_name`            varchar(255),
  `attack_severity`        enum('NONE','LOW','MEDIUM','HIGH','UNKNOWN') NOT NULL,
  `start_time`             int(11),
  `end_time`               int(11),
  `created`                datetime DEFAULT NULL,
  `updated`                datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_attack_detail trigger when any attribute of uri_filtering_attack_detail change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_attack_detail_trigger AFTER UPDATE ON uri_filtering_attack_detail
FOR EACH ROW
BEGIN
  IF NEW.vendor_id <> OLD.vendor_id OR NEW.attack_id <> OLD.attack_id OR NEW.attack_name <> OLD.attack_name OR NEW.attack_severity <> OLD.attack_severity
     OR NEW.start_time <> OLD.start_time OR NEW.end_time <> OLD.end_time THEN
    SELECT MySQLNotification('uri_filtering_attack_detail', NEW.tele_pre_mitigation_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_source_count
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_source_count`;

CREATE TABLE `uri_filtering_source_count` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `low_percentile_g`      int(11),
  `mid_percentile_g`      int(11),
  `high_percentile_g`     int(11),
  `peak_g`                int(11),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


# uri_filtering_source_count trigger when any attribute of uri_filtering_source_count change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_source_count_trigger AFTER UPDATE ON uri_filtering_source_count
FOR EACH ROW
BEGIN
  IF NEW.low_percentile_g <> OLD.low_percentile_g OR NEW.mid_percentile_g <> OLD.mid_percentile_g OR NEW.high_percentile_g <> OLD.high_percentile_g OR NEW.peak_g <> OLD.peak_g THEN
    SELECT MySQLNotification('uri_filtering_source_count', NEW.tele_attack_detail_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_top_talker
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_top_talker`;

CREATE TABLE `uri_filtering_top_talker` (
  `id`                    bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_attack_detail_id` bigint(20) NOT NULL,
  `spoofed_status`        tinyint(1),
  `created`               datetime DEFAULT NULL,
  `updated`               datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_top_talker trigger when any attribute of uri_filtering_top_talker change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_top_talker_trigger AFTER UPDATE ON uri_filtering_top_talker
FOR EACH ROW
BEGIN
  IF NEW.spoofed_status <> OLD.spoofed_status THEN
    SELECT MySQLNotification('uri_filtering_top_talker', NEW.tele_attack_detail_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_source_prefix
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_source_prefix`;

CREATE TABLE `uri_filtering_source_prefix` (
  `id`                 bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20)   NOT NULL,
  `addr`               varchar(255) DEFAULT NULL,
  `prefix_len`         int(11)      DEFAULT NULL,
  `created`            datetime     DEFAULT NULL,
  `updated`            datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

# uri_filtering_source_prefix trigger when any attribute of uri_filtering_source_prefix change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_source_prefix_trigger AFTER UPDATE ON uri_filtering_source_prefix
FOR EACH ROW
BEGIN
  IF NEW.addr <> OLD.addr OR NEW.prefix_len <> OLD.prefix_len THEN
    SELECT MySQLNotification('uri_filtering_source_prefix', NEW.tele_top_talker_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_source_port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_source_port_range`;

CREATE TABLE `uri_filtering_source_port_range` (
  `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20) NOT NULL,
  `lower_port`         int(11)    NOT NULL,
  `upper_port`         int(11)    DEFAULT NULL,
  `created`            datetime   DEFAULT NULL,
  `updated`            datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_source_port_range trigger when any attribute of uri_filtering_source_port_range change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_source_port_range_trigger AFTER UPDATE ON uri_filtering_source_port_range
FOR EACH ROW
BEGIN
  IF NEW.lower_port <> OLD.lower_port OR NEW.upper_port <> OLD.upper_port THEN
    SELECT MySQLNotification('uri_filtering_source_port_range', NEW.tele_top_talker_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# uri_filtering_icmp_type_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `uri_filtering_icmp_type_range`;

CREATE TABLE `uri_filtering_icmp_type_range` (
  `id`                 bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_top_talker_id` bigint(20) NOT NULL,
  `lower_type`         int(11)    NOT NULL,
  `upper_type`         int(11)    DEFAULT NULL,
  `created`            datetime   DEFAULT NULL,
  `updated`            datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# uri_filtering_icmp_type_range trigger when any attribute of uri_filtering_icmp_type_range change
# ------------------------------------------------------------------------------

DELIMITER @@
CREATE TRIGGER uri_filtering_icmp_type_range_trigger AFTER UPDATE ON uri_filtering_icmp_type_range
FOR EACH ROW
BEGIN
  IF NEW.lower_type <> OLD.lower_type OR NEW.upper_type <> OLD.upper_type THEN
    SELECT MySQLNotification('uri_filtering_icmp_type_range', NEW.tele_top_talker_id) INTO @x;
  END IF;
END@@
DELIMITER ;

# vendor_mapping
# ------------------------------------------------------------

DROP TABLE IF EXISTS `vendor_mapping`;

CREATE TABLE `vendor_mapping` (
  `id`             bigint(20) NOT NULL AUTO_INCREMENT,
  `data_client_id` bigint(20) NOT NULL,
  `vendor_id`      int(11)    NOT NULL,
  `created`        datetime   DEFAULT NULL,
  `updated`        datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# attack_mapping
# ------------------------------------------------------------

DROP TABLE IF EXISTS `attack_mapping`;

CREATE TABLE `attack_mapping` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `vendor_mapping_id` bigint(20)   NOT NULL,
  `attack_id`         int(11)      NOT NULL,
  `attack_name`       varchar(255) NOT NULL,
  `created`           datetime     DEFAULT NULL,
  `updated`           datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;