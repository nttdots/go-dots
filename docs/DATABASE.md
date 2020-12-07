# Setting up databases

In this document, we describe how to configure the data in the database.

## Editing the data on the database 

For the current implementation, some part of the system configuration is done by editing the data on the database directly. Edit dots_server/db_models/template.sql to fit in your environment. Those configurations are blocker configurations and customer configurations. We detail each configration in the following sections.

After completing the edit, you can update the data in the database by the command below.

mysql -u root dots < dots_server/db_models/template.sql 

## Blocker configuration

The blockers are the entities to mitigate the DDoS attacks based on the requests from the dots_server. The blocker related tables are below:

* 'blocker': Basic blocker information of each blocker_type
* 'blocker_parameter': Additional configuration for the blockers
* 'blocker_configuration': The configuration of the blockers for each customer.
* 'blocker_configuration_parameter': Additional configuration for the blockers corresponding to different customers
* 'login_profile': The login profile for the blockers. The current implementation does not utilize these profile.

The tables you have to edit here are 'blocker', 'blocker_parameter', 'blocker_configuration' and 'blocker_configuration_parameter'.

### 'blocker' table

The current implementation support three blocker types: 'GoBGP-RTBH', 'Arista-ACL' and 'GoBGP-FlowSpec'. So set the type field of every blockers to 'GoBGP-RTBH', 'Arista-ACL' or 'GoBGP-FlowSpec'. In the capacity field, you can specify the strength of the blocker. The capacity is the number of the entries the blocker can deal with. The load field is the number of the entries the blocker is now dealing with and will be updated by the system. You should set to 0 for the initial configuration.

### 'blocker_parameter' table

The 'blocker_parameter' tables is for the further configurations of each blocker. The 'blocker_parameter' has the key/value structure. The parameter has 3 keys below:

* nextHop: the nextHop IP address for the RTBH configuration. If you configure the nextHop as '0.0.0.0', the blocker simply blackhole the attack traffics.
* host: IP address or FQDN of the blocker.
* port: port number of the control access.

### 'blocker_configuration' table

The 'blocker_configuration' tables is for the blocker configurations for each blocker_type. There are three fields that can be modified:

* customer_id: the certificate id of the DOTS client that are being serviced by DOTS server
* target_type: the type of targets which are applied to the blocker, there are two types: 'mitigation_request' and 'datachannel_acl'
* blocker_type: the type of blockers, there are three types: 'GoBGP-RTBH', 'Arista-ACL' and 'GoBGP-FlowSpec'

The blocker_configuration data should be provided to both 'mitigation_request' and 'datachannel_acl' for each customer. For each 'target_type', we can set 'blocker_type' value follow as:

* 'mitigation_request': 'GoBGP-RTBH', 'Arista-ACL' or 'GoBGP-FlowSpec'
* 'datachannel_acl': 'Arista-ACL' and 'GoBGP-FlowSpec'

### 'blocker_configuration_parameter' table

The 'blocker_configuration_parameter' tables is for the further configurations of each blocker corresponding to different customers. The 'blocker_configuration_parameter' has the key/value structure. The parameter has 3 keys below:

* vrf: the virtual routing forwarding of GoBGP Flowspec that used for redirecting (Ex: 1.1.1.1:100). This field is only used when blocker_type is 'GoBGP-FlowSpec' and target-type is 'mitigation_request'
* arista_connection: the name of Arista configuration in ~/.eapi.conf file (Ex: arista). This field is only used when blocker_type is 'Arista-ACL'
* arista_interface: the name of Arista interface that to apply 'Arista-ACL' (Ex: Ethernet 1). This field is only used when blocker_type is 'Arista-ACL'

## Customer configuration

Our system utilizes the concept of 'customer' to validate (mainly mitigation) requests as described in the "How it works" section of [README](../README.md).  The customer related tables are listed below.

* 'customer': Basic blocker information
* 'customer_common_name': customer certificate common names
* 'parameter_value': message parameters which can be utilized for the message validations
* 'prefix': IP address prefix information which can be utilized for the message validations

We detail each table in the following sections.

### 'customer'

There is no special instruction for this table.

### 'customer_common_name' table

 To validate each messages, first we check the common name field of the client certificates. To validate the common names, we utilize this table.

### 'parameter_value'

This tables is to store the identifier parameters and mitigation scope parameters. If you specify the customer_id field appropriately, the entry is utilized for the mitigation request validation for the customer specified by the customer_id.

The types of parameters for the validation use are below.

* FQDN: valid FQDN(s) for the customer
* URI: valid URI(s) for the customer

### 'prefix'

The 'prefix' table is to store the IP address information. If the mitigation scopes in the mitigation requests includes the IP address prefixes, the message validator validates the messages by checking whether the database has the entries which has the 'type' field is 'ADDRESS_RANGE' and the CIDR block specified by the 'addr' field and the 'prefix_len' matches to the IP address block specified by the mitigation scope of the message.

## Vendor attack mapping

Vendor attack mapping includes the client's vendor mapping and the server's vendor mapping.

* The client's vendor mapping data is sent from dots client.
* The server's vendor mapping data: we manual insert data in mysql as follows:

    INSERT INTO `vendor_mapping` (`id`, `data_client_id`, `vendor_id`, `vendor_name`, `last_updated`,`created`, `updated`)
    VALUES (1, 0, 345, 'mitigator-c', 1576812345, '2020-05-18 16:44:34', '2020-05-18 16:44:34');

    INSERT INTO `attack_mapping` (`id`, `vendor_mapping_id`, `attack_id`, `attack_description`,`created`, `updated`)
    VALUES (1, 1, 1, 'attack-description 1', '2020-05-18 16:44:34', '2020-05-18 16:44:34'),
           (2, 1, 2, 'attack-description 2', '2020-05-18 16:44:34', '2020-05-18 16:44:34');

