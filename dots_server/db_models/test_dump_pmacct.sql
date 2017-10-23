# acct_v5
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acct_v5`;

CREATE TABLE acct_v5 (
  `agent_id` INT(4) UNSIGNED NOT NULL,
  `class_id` CHAR(16) NOT NULL,
  `mac_src` CHAR(17) NOT NULL,
  `mac_dst` CHAR(17) NOT NULL,
  `vlan` INT(2) UNSIGNED NOT NULL,
  `ip_src` CHAR(15) NOT NULL,
  `ip_dst` CHAR(15) NOT NULL,
  `src_port` INT(2) UNSIGNED NOT NULL,
  `dst_port` INT(2) UNSIGNED NOT NULL,
  `ip_proto` CHAR(6) NOT NULL,
  `tos` INT(4) UNSIGNED NOT NULL,
  `packets` INT UNSIGNED NOT NULL,
  `bytes` BIGINT UNSIGNED NOT NULL,
  `flows` INT UNSIGNED NOT NULL,
  `stamp_inserted` DATETIME NOT NULL,
  `stamp_updated` DATETIME,
  PRIMARY KEY (`agent_id`, `class_id`, `mac_src`, `mac_dst`, `vlan`, `ip_src`, `ip_dst`, `src_port`, `dst_port`, `ip_proto`, `tos`, `stamp_inserted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `acct_v5` (`agent_id`, `class_id`, `mac_src`, `mac_dst`, `vlan`, `ip_src`, `ip_dst`, `src_port`, `dst_port`, `ip_proto`, `tos`, `packets`, `bytes`, `flows`, `stamp_inserted`, `stamp_updated`)
VALUES
  (1,1,'A0:B2:D5:7F:81:B3','A0:B2:D6:1F:15:22',10,'101.123.200.211','192.168.1.1',8989,8990,'123456',100,820,24543,1,NOW(),NOW()),
  (1,1,'A0:B2:D5:7F:81:B3','A0:B2:D6:1F:15:22',10,'101.123.200.211','192.168.1.1',8989,8990,'123456',100,100,12135,1,NOW()+INTERVAL 1 MINUTE,NOW()+INTERVAL 1 MINUTE),
  (1,1,'A0:B2:D5:7F:81:B3','A0:B2:D6:1F:15:22',10,'101.123.200.211','192.168.1.1',8989,8990,'123456',100,1045,54439,1,NOW()+INTERVAL 2 MINUTE,NOW()+INTERVAL 2 MINUTE),
  (2,2,'B0:B2:D5:7F:81:B3','B0:B2:D6:1F:15:22',10,'201.123.200.212','192.168.1.101',8989,8990,'4444',100,20,2135,0,NOW(),NOW()),
  (2,2,'B0:B2:D5:7F:81:B3','B0:B2:D6:1F:15:22',10,'201.123.200.212','192.168.1.101',8989,8990,'4444',100,12354,7464835,0,NOW()+INTERVAL 1 MINUTE,NOW()+INTERVAL 1 MINUTE),
  (2,2,'B0:B2:D5:7F:81:B3','B0:B2:D6:1F:15:22',10,'201.123.200.212','192.168.1.101',8989,8990,'4444',100,8584,348304,0,NOW()+INTERVAL 2 MINUTE,NOW()+INTERVAL 2 MINUTE),
  (2,2,'B0:B2:D5:7F:81:B3','B0:B2:D6:1F:15:22',10,'201.123.200.212','192.168.1.101',8989,8990,'4444',100,494,439433,0,NOW()+INTERVAL 3 MINUTE,NOW()+INTERVAL 3 MINUTE);
