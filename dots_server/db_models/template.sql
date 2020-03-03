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

####### Add an insert query to add the basic data of your blockers here.
#  id: id is the identifier of the blocker object
#  blocker_type: currently only 'GoBGP-RTBH' is supported
#  capacity: capacity with which this blocker can deal with the attack traffics
#  load: load of the traffic this blocker is currently dealing with.
#
# example query:
#  INSERT INTO `blocker` (`id`, `blocker_type`, `capacity`, `load`, `created`, `updated`)
#  VALUES
#   (1,'GoBGP-RTBH', 100, 0,'2017-04-13 13:44:34','2017-04-13 13:44:34')
#   (2,'Arista-ACL', 100, 0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#   (3,'GoBGP-FlowSpec', 100, 0,'2017-04-13 13:44:34','2017-04-13 13:44:34');


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

####### Add an insert query to add detailed information about the blockers you inserted above
#  id: id of this blocker_parameter. note that this id is not the blocker's id.
#  blocker_id: id of the blocker this parameter corresponds to.
#  key: the parameter type of this blocker parameter. the values are 'nextHop', 'host' and 'port'.
#    nextHop: nextHop IP address of the DDoS traffic
#    host: the blocker's FQDN or IP address
#    port: the API port if exists
#  value: value for the key
#
# example query:
#  INSERT INTO `blocker_parameter` (`id`, `blocker_id`, `key`, `value`, `created`, `updated`)
#  VALUES
#  (1, 1, 'nextHop', '0.0.0.1','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (2, 1, 'host', '127.0.0.1','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (3, 1, 'port', '50051','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (4, 2, 'nextHop', '0.0.0.1','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (5, 2, 'host', '127.0.0.1','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (6, 2, 'port', '50051','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (7, 3, 'nextHop', '0.0.0.2','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (8, 3, 'host', '127.0.0.1','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (9, 3, 'port', '50051','2017-04-13 13:44:34','2017-04-13 13:44:34');


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

####### Add an insert query to add the basic data of the customers here.
#  id: id is the identifier of the customer object
#  common_name: name of the customer. this is for you to distinguish each customer easily.
#
# example query:
#  INSERT INTO `customer` (`id`, `common_name`, `created`, `updated`)
#  VALUES
#  (1,'name','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (2,'localhost','2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (3,'local-host','2017-04-13 13:44:34','2017-04-13 13:44:34');

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

####### Add an insert query to add customer related parameters.
####### these parameters are to store the identifier parameters or mitigation scope parameters.
#  id: id of this login profile.
#  customer_id: id of the customer this parameter corresponds to. If specified, this parameter object is used for the validation use.
#  mitigation_scope_id: id of a mitigation scope the customer requested  basically set by the system. if this parameter object is not for a mitigation scope, set this field to '0'. Note that this id is not the 'mitigation_id', but the database id.
#  type: type of this parameter. these are based on the internet drafts.
#   'TARGET_PROTOCOL','FQDN','URI','TRAFFIC_PROTOCOL','ALIAS_NAME'
#  string_value: if this parameter is a type of string parameter, specify the value in the string format.
#  int_value: if this parameter is a type of int parameter, specify the value in the integer format.
#
# example query:
#  INSERT INTO `parameter_value` (`id`, `customer_id`, `mitigation_scope_id`, `type`, `string_value`, `int_value`, `created`, `updated`)
#  VALUES
#  (1,1,0,'FQDN','golang.org',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (3,2,0,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (4,3,0,'FQDN','local-host',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (5,0,1,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#  (6,0,2,'FQDN','client.sample.example.com',0,'2017-04-13 13:44:34','2017-04-13 13:44:34');


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
####### Basically the table 'port_range' is modified by the system only.

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

####### Add an insert query to add customer related parameters.
#  id: id of this login profile.
#  customer_id: id of the customer this prefix corresponds to. If specified, this prefix object is used for the validation use.
#  mitigation_scope_id: id of a mitigation scope the customer requested. basically set by the system. if this prefix object is not for a mitigation scope, set this field to '0'. Note that this id is not the 'mitigation_id', but the database id.
#  type: type of this parameter. part of these are based on drafts of the data channel and the signal channel.
#    'TARGET_PREFIX','SOURCE_IPV4_NETWORK','DESTINATION_IPV4_NETWORK','IP','PREFIX','ADDRESS_RANGE','IP_ADDRESS','TARGET_IP'
#  addr: address of the prefix
#  prefix_len: length of the prefix
#
# example query:
#  INSERT INTO `prefix` (`id`, `customer_id`, `mitigation_scope_id`, `type`, `addr`, `prefix_len`, `created`, `updated`)
#  VALUES
#    (1,1,0,'ADDRESS_RANGE','192.168.1.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (5,2,0,'ADDRESS_RANGE','127.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (6,2,0,'ADDRESS_RANGE','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (7,2,0,'ADDRESS_RANGE','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (8,2,0,'ADDRESS_RANGE','192.168.7.0',24,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (9,3,0,'ADDRESS_RANGE','129.0.0.1',32,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (10,3,0,'ADDRESS_RANGE','2003:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (11,3,0,'ADDRESS_RANGE','2003:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (12,0,1,'TARGET_IP','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (13,0,1,'TARGET_PREFIX','2002:db8:6401::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (14,0,2,'TARGET_IP','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),
#    (15,0,2,'TARGET_PREFIX','2002:db8:6402::',64,'2017-04-13 13:44:34','2017-04-13 13:44:34'),


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

####### Basically the table 'mitigation_scope' is modified by the system only.

# mitigation_scope trigger when status change
# ------------------------------------------------------------

DROP FUNCTION IF EXISTS MySQLNotification;
CREATE FUNCTION MySQLNotification RETURNS INTEGER SONAME 'mysql-notification.so';

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

####### Basically the table 'signal_session_configuration' is modified by the system only.

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

####### Basically the table 'protection' is modified by the system only.

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

####### Basically the table 'go_bgp_parameter' is modified by the system only.

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

####### Basically the table 'protection_status' is modified by the system only.


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

####### Basically the table 'data_clients' is modified by the system only.

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

####### Basically the table 'data_aliases' is modified by the system only.

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

####### Basically the table 'data_acls' is modified by the system only.

# arista_parameter
# ------------------------------------------------------------

DROP TABLE IF EXISTS `arista_parameter`;

CREATE TABLE `arista_parameter` (
  `id`                  bigint(20) NOT NULL AUTO_INCREMENT,
  `protection_id`       bigint(20) NOT NULL,
  `acl_type`            varchar(255) NOT NULL,
  `acl_filtering_rule`  text     NOT NULL,
  `created`             datetime DEFAULT NULL,
  `updated`             datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'arista_parameter' is modified by the system only.

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

####### Add an insert query to add detailed information about the configuration of blockers you inserted above for each customer
#  id: id of this blocker_configuration. note that this id is not the blocker's id.
#  customer_id: id of the customer that is configured
#  target_type: the type of target. the values are 'mitigation_request' or 'datachannel_acl'.
#  blocker_type: the type of blocker that is used for protecting the target above. the values are 'GoBGP-RTBH', 'Arista-ACL' or 'GoBGP-FlowSpec'
#
# example query:
#  INSERT INTO `blocker_configuration` (`id`, `customer_id`, `target_type`, `blocker_type`, `created`, `updated`)
#  VALUES
#  (1, 1, "mitigation_request", "GoBGP-RTBH", '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
#  (2, 1, "datachannel_acl", "Arista-ACL", '2017-04-13 13:44:34', '2017-04-13 13:44:34');

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

####### Add an insert query to add detailed information about the configuration of blockers you inserted above
#  id: id of this blocker_configuration_parameter. note that this id is not the blocker_configuration's id.
#  blocker_configuration_id: id of the blocker configuration this parameter corresponds to.
#  key: the parameter type of this blocker parameter. the values are 'vrf', 'aristaConnection' and 'aristaInterface'.
#    vrf: virtual routing forwarding address of the GoBGp-FlowSpec blocker
#    aristaConnection: the connection name that is used for Arista blocker
#    aristaInterface: the interface name that is used for registering Arista ACL
#  value: value for the key
#
# example query:
#  INSERT INTO `blocker_configuration_parameter` (`id`, `blocker_configuration_id`, `key`, `value`, `created`, `updated`)
#  VALUES
#  (1, 1, 'vrf', '1.1.1.1:100', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
#  (2, 1, 'aristaConnection', 'arista', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
#  (3, 1, 'aristaInterface', 'Ethernet 1', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
#  (4, 2, 'aristaConnection', 'arista', '2017-04-13 13:44:34', '2017-04-13 13:44:34'),
#  (5, 2, 'aristaInterface', 'Ethernet 1', '2017-04-13 13:44:34', '2017-04-13 13:44:34');


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

####### Basically the table 'flow_spec_parameter' is modified by the system only.


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

####### Basically the table 'control_filtering' is modified by the system only.


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

####### Basically the table 'telemetry_setup' is modified by the system only.

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

####### Basically the table 'telemetry_configuration' is modified by the system only.

# unit_configuration
# ------------------------------------------------------------

DROP TABLE IF EXISTS `unit_configuration`;

CREATE TABLE `unit_configuration` (
  `id`             bigint(20) NOT NULL AUTO_INCREMENT,
  `tele_config_id` bigint(20) NOT NULL,
  `unit`           enum('PPS','KILO_PPS','BPS','KILOBYTES_PS','MEGABYTES_PS','GIGABYTES_PS') NOT NULL,
  `unit_status`    tinyint(1) DEFAULT NULL,
  `created`        datetime   DEFAULT NULL,
  `updated`        datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'unit_configuration' is modified by the system only.

# total_pipe_capability
# ------------------------------------------------------------

DROP TABLE IF EXISTS `total_pipe_capacity`;

CREATE TABLE `total_pipe_capacity` (
  `id`            bigint(20)   NOT NULL AUTO_INCREMENT,
  `tele_setup_id` bigint(20)   NOT NULL,
  `link_id`       varchar(255) DEFAULT NULL,
  `capacity`      int(11)      DEFAULT NULL,
  `unit`          enum('PPS','KILO_PPS','BPS','KILOBYTES_PS','MEGABYTES_PS','GIGABYTES_PS') NOT NULL,
  `created`       datetime     DEFAULT NULL,
  `updated`       datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'total_pipe_capacity' is modified by the system only.

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

####### Basically the table 'baseline' is modified by the system only.

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

####### Basically the table 'telemetry_prefix' is modified by the system only.

# telemetry_port_range
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_port_range`;

CREATE TABLE `telemetry_port_range` (
  `id`         bigint(20) NOT NULL AUTO_INCREMENT,
  `type`       enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`    bigint(20) NOT NULL,
  `lower_port` int(11)    NOT NULL,
  `upper_port` int(11)    DEFAULT NULL,
  `created`    datetime   DEFAULT NULL,
  `updated`    datetime   DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'telemetry_port_range' is modified by the system only.

# telemetry_parameter_value
# ------------------------------------------------------------

DROP TABLE IF EXISTS `telemetry_parameter_value`;

CREATE TABLE `telemetry_parameter_value` (
  `id`             bigint(20)   NOT NULL AUTO_INCREMENT,
  `type`           enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`        bigint(20)   NOT NULL,
  `parameter_type` enum('TARGET_PROTOCOL','FQDN','URI') NOT NULL,
  `string_value`   varchar(255) DEFAULT NULL,
  `int_value`      int(11)      DEFAULT NULL,
  `created`        datetime     DEFAULT NULL,
  `updated`        datetime     DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'telemetry_parameter_value' is modified by the system only.

# traffic
# ------------------------------------------------------------

DROP TABLE IF EXISTS `traffic`;

CREATE TABLE `traffic` (
  `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
  `customer_id`       int(11)      NOT NULL,
  `cuid`              varchar(255) NOT NULL,
  `type`              enum('TELEMETRY','TELEMETRY_SETUP') NOT NULL,
  `type_id`           bigint(20)   NOT NULL,
  `traffic_type`      enum('TOTAL_TRAFFIC_NORMAL_BASELINE','TOTAL_ATTACK_TRAFFIC','TOTAL_TRAFFIC') NOT NULL,
  `unit`              enum('PPS','KILO_PPS','BPS','KILOBYTES_PS','MEGABYTES_PS','GIGABYTES_PS') NOT NULL,
  `protocol`          int(11)     NOT NULL,
  `low_percentile_g`  int(11)     DEFAULT NULL,
  `mid_percentile_g`  int(11)     DEFAULT NULL,
  `high_percentile_g` int(11)     DEFAULT NULL,
  `peak_g`            int(11)     DEFAULT NULL,
  `created`           datetime    DEFAULT NULL,
  `updated`           datetime    DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####### Basically the table 'traffic' is modified by the system only.

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

####### Basically the table 'total_connection_capacity' is modified by the system only.