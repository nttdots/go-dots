# Setting up databases

In this document, we describe how to configure the data in the database.

## Editing the data on the database 

For the current implementation, some part of the system configuration is done by editing the data on the database directly. Edit dots_server/db_models/template.sql to fit in your environment. Those configurations are blocker configurations and customer configurations. We detail each configration in the following sections.

After completing the edit, you can update the data in the database by the command below.

mysql -u root dots < dots_server/db_models/template.sql 

## Blocker configuration

The blockers are the entities to mitigate the DDoS attacks based on the requests from the dots_server. The blocker related tables are below:

* 'blocker': Basic blocker information
* 'blocker_parameter': Additinal configuration for the blockers
* 'login_profile': The login profile for the blocker. The current implementation does not utilize these profile.

The tables you have to edit here are 'blocker' and 'blocker_parameter'.

### 'blocker' table

The current implementation only support the blocker type of 'GoBGP-RTBH'. So set the type field of every blockers to 'GoBGP-RTBH'. In the capacity field, you can specify the strength of the blocker. The capacity is the number of the entries the blocker can deal with. The load field is the number of the entries the blocker is now dealing with and will be updated by the system. You should set to 0 for the initial configuration.

### 'blocker_parameter' table

The 'blocker_parameter' tables is for the further configurations of each blocker. 'blocker_parameter' has the key/value structure. The parameter has 3 keys below:

* nextHop: the nextHop IP address for the RTBH configuration. If you configure the nextHop as '0.0.0.0', the blocker simply blackhole the attack traffics.
* host: IP address or FQDN of the blocker.
* port: port number of the controll access.

## Customer configuration

Our system utilizes the concept of 'customer' to validate (mainly mitigation) requests as described in the "How it works" section of [README](../README.md).  The customer related tables are listed below.

* 'customer': Basic blocker information
* 'customer_common_name': customer certificate common names
* 'customer_radius_user': radius認証に使用するユーザー/パスワード
* 'parameter_value': message parameters which can be utilized for the message validations
* 'prefix': IP address prefix information which can be utilized for the message validations

We detail each table in the following sections.

### 'customer'

There is no special instruction for this table.

### 'customer_common_name' table

 To validate each messages, first we check the common name field of the client certificates. To validate the common names, we utilize this table.

### 'customer_radius_user' table

 メッセージの認証にradiusを使用する場合、このテーブルのuser_name, user_passwordフィールドの内容を使用します。
 user_realmに文字が設定されていた場合、radiusサーバーに送信するユーザー名は、\[user_name\]@\[user_password\] となります。

### 'parameter_value'

This tables is to store the identifier parameters and mitigation scope parameters. If you specify the customer_id field appropriately, the entry is utilized for the mitigation request validation for the customer specified by the customer_id.

The types of parameters for the validation use are below.

* FQDN: valid FQDN(s) for the customer
* URI: valid URI(s) for the customer
* E_164: valid E.164 nubmers for the cusotmer

### 'prefix'

The 'prefix' table is to store the IP address information. If the mitigation scopes in the mitigation requests includes the IP address prefixes, the message validator validates the messages by checking whether the database has the entries which has the 'type' field is 'ADDRESS_RANGE' and the CIDR block specified by the 'addr' field and the 'prefix_len' matches to the IP address block specified by the mitigation scope of the message.

