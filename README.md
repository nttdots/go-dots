

# go-dots

![logo](https://github.com/nttdots/go-dots/blob/master/go-dots_logo/go-dots_logo_blue.png)

"go-dots" is a DDoS Open Threat Signaling (dots) implementation written in Go. This implmentation is based on the Internet drafts below. 

* draft-ietf-dots-signal-channel-18
* draft-ietf-dots-data-channel-02 
* draft-ietf-dots-architecture-05 
* draft-ietf-dots-requirements-11 
* draft-ietf-dots-use-cases-09 

This implementation is not fully compliant with the documents listed above.  For example, we are utilizing CoAP as the data channel protocol while the current version of the data channel document specifies RESTCONF as the data channel protocol.

Licensed under Apache License 2.0.

# How to Install

## Requirements

* gcc 5.4.0 or higher
* make, autoconf, automake, libtool, pkg-config, pkgconf or pkg-config
* [git](https://git-scm.com/)
* [go](https://golang.org/doc/install)
  * go 1.9 or later is required. (for latest GoBGP)
  * set PATH to go and set $GOPATH, using their instructions.
* [openssl](https://www.openssl.org/)
  * OpenSSL 1.1.0 or higher (for libcoap)

* MySQL 5.7 or higher and its development package
  * Install mysql development package in Ubuntu:
    $ sudo apt-get install libmysqld-dev

## Recommandation Environment
* Ubuntu 16.04+
* macOS High Sierra 10.13+

## How to build go-dots
### Build libcoap for go-dots

Currenly supported libcoap version : 1365dea

    $ git clone https://github.com/obgm/libcoap.git
    $ cd libcoap
    $ git checkout 1365dea39a6129a9b7e8c579537e12ffef1558f6
    $ ./autogen.sh
    $ ./configure --disable-documentation --with-openssl
    $ make
    $ sudo make install
    
### Install go-dots
To install go-dots source codes and command line programs, use the following command:
    
    $ go get -u github.com/nttdots/go-dots/...
    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ make && make install

# How to install (In Japanese)

* [qiita](https://qiita.com/__kaname__/items/774819444de2a12d99b9)

# How it works

go-dots consists of 3 daemons below:

* dots_server: dots server agent
* dots_client: dots client specified in the Internet drafts
* dots_client_controller: The controller for the dots_client. Operators controll the dots_client by this controller

To explain the relationships between each daemons in detail, here we illustrates the mitigation request procedure of our system. First, dots_client_controller passes request JSONs to the dots_client. Examples of these JSONs are located in the 'dots_client/' directory in this Github project. Remember to edit the 'mitigation-id' fields in the JSON file every time you make requests to the dots_server. Then the dots client converts the received JSON to the CBOR format and send a mitigation request to the server via a CoAP connection.

The figure below shows the detailed sequence diagram which depicts how a dots_server handles a mitigation request. 

<img src='https://github.com/nttdots/go-dots/blob/documentation/docs/pics/mitigation_request_sequence.png' title='mitigation requests sequence diagram'>

Upon receiving CoAP messages, a dots_server first find the appropriate message controller for the received message. The message controller first find the customer information bound with the common name field contained in the client certificate. The customer information is configured and stored in the RDB before the server is started. To setup databases, refer to [Setting up databases](./docs/DATABASE.md). If no appropriate customer objets is found, the dots_server decline the received request. 

Then the server retrieves mitigation scopes contained in the received mitigation request message. Again the dots_server validate the mitigation scopes to determine whether the customer which issued the request has the valid privilege for the mitigation operations. If the validation is successfully completed, the server select a blocker for the mitigation and execute the mitigation.

# Usage

## Server Configuration

Server Configuration is done by the system configuration file and the database setup. The system configuration file is specified via '-config' option when the 'dots_server' is invoked. The sample configuration files are located as 'dots_server/dots_server.yaml' and 'dots_server/dots_server.yaml.template'. 

To set up your database, refer to the [Database configuration document](./docs/DATABASE.md)

## Server
    $ $GOPATH/bin/dots_server --config [config.yml file (ex: go-dots/dots_server/dots_server.yaml)]


## Client
    $ $GOPATH/bin/dots_client --server localhost --signalChannelPort=5684 --config [config.yml file (ex: go-dots/dots_client/dots_client.yaml)] -vv

## GoBGP Server
To install and run gobgp-server, refer to the following link:
* [gobgp-server](https://github.com/osrg/gobgp)

    
### Client Controller [mitigation_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraft.json
   
### Client Controller [mitigation_retrieve_all]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw

### Client Controller [mitigation_retrieve_one]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123

### Client Controller [mitigation_withdraw]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Delete \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123

### Client Controller [mitigation_observe]
A DOTS client can convey the 'observe' option set to '0' in the GET request to receive notification whenever status of mitigation request changed
and unsubscribe itself by issuing GET request with 'observe' option set to '1'

Subscribe for resource observation:

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -observe=0

Unsubscribe from resource observation:

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -observe=1

Subscriptions are valid as long as current session exists. When session is renewed (e.g DOTS client does not receive response from DOTS server for its Ping message 
in a period of time, it decided that server has been disconnected, then re-connects), all previous subscriptions will be lost. In such cases, DOTS clients will have to re-subscribe
for observation. Below is recommended step: 

    ・GET a list of all existing mitigations (that were created before server restarted)
    ・PUT mitigations  one by one
    ・GET + Observe for all the mitigations that should be observed

### Client Controller [mitigation_efficacy_update]
A DOTS client can convey the 'If-Match' option with empty value in the PUT request to transmit DOTS mitigation efficacy update to the DOTS server:
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -ifMatch="" \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraftEfficacyUpdate.json

### Client Controller [session_configuration_request]
    $ $GOPATH/bin/dots_client_controller -request session_configuration -method Put \
     -sid 234 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleSessionConfigurationDraft.json

### Client Controller [session_configuration_retrieve_default]
    $ $GOPATH/bin/dots_client_controller -request session_configuration -method Get

### Client Controller [session_configuration_retrieve_one]
    $ $GOPATH/bin/dots_client_controller -request session_configuration -method Get \
      -sid 234

### Client Controller [session_configuration_delete]
    $ $GOPATH/bin/dots_client_controller -request session_configuration -method Delete \
      -sid 234

###  Client Controller [client_configuration_request]
DOTS signal channel session configuration supports 2 sets of parameters : 'mitigating-config' and 'idle-config'.
The same or distinct configuration set may be used during times when a mitigation is active ('mitigating-config') and when no mitigation is active ('idle-config').
Dots_client uses 'idle-config' parameter set by default. It can be configured to switch to the other parameter set by client_configuration request

Configure dots_client to use 'idle-config' parameters

    $ $GOPATH/bin/dots_client_controller -request client_configuration -method POST \
    -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleClientConfigurationRequest_Idle.json

Configure dots_client to use 'mitigating-config' parameters

    $ $GOPATH/bin/dots_client_controller -request client_configuration -method POST \
    -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleClientConfigurationRequest_Mitigating.json

##  Data Channel

All shell-script and sample json files are located in below directory:
    $ cd $GOPATH/src/github.com/nttdots/go-dots/dots_client/data/

### Get Root Resource Path

    Get root resource:
    $ ./get_root_resource.sh SERVER_NAME

    Example:
        - Request: $ ./get_root_resource.sh https://127.0.0.1:10443
        - Response:
        <XRD xmlns="https://127.0.0.1">
            <Link rel="restconf" href="https://127.0.0.1:10443/v1/restconf"></Link>
        </XRD>

    "{href}" value will be used as the initial part of the path in the request URI of subsequent requests

### Managing DOTS Clients
Registering DOTS Clients

    Post dots_client:
    $ ./do_request_from_file.sh POST {href}/data/ietf-dots-data-channel:dots-data sampleClient.json

    Put dots_client:
    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleClient.json

Uregistering DOTS Clients

    $ ./do_request_from_file.sh DELETE {href}/data/ietf-dots-data-channel:dots-data/dots-client=123

### Managing DOTS Aliases
Create Aliases

    Post alias:
    $ ./do_request_from_file.sh POST {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleAlias.jon

    Put alias:
    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleAlias.jon

Retrieve Installed Aliases

    Get all aliases without 'content' parameter (default is get all type attributes, including configurable and non-configurable attributes):
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases

    Get all aliases with 'content'='config' (get configurable attributes only):
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases?content=config

    Get all aliases with 'content'='nonconfig' (get non-configurable attributes only):
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases?content=nonconfig

    Get all aliases with 'content'='all'(get all type attributes, including configurable and non-configurable attributes):
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases?content=all

    Get specific alias without 'content' parameter:
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=https1

    Get specific alias with 'content'='config':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=https1?content=config

    Get specific alias with 'content'='nonconfig':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=https1?content=nonconfig

    Get specific alias with 'content'='all':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=https1?content=all

Delete Aliases

    $ ./do_request_from_file.sh DELETE {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=https1

### Managing DOTS Filtering Rules
Retrieve DOTS Filtering Capabilities

    Get Capabilities without 'content' parameter:
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/capabilities

    Get Capabilities with 'content'='config':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/capabilities?content=config

    Get Capabilities with 'content'='nonconfig':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/capabilities?content=nonconfig

    Get Capabilities with 'content'='all':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/capabilities?content=all

Install Filtering Rules

    Post acl:
    $ ./do_request_from_file.sh POST {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleAcl.json

    Put acl:
    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleAcl.json

Retrieve Installed Filtering Rules

    Get all Acl without 'content' parameter:
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls

    Get all Acl with 'content'='config':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls?content=config

    Get all Acl with 'content'='nonconfig':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls?content=nonconfig

    Get all Acl with 'content'='all':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls?content=all

    Get specific acl without 'content' parameter:
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=sample-ipv4-acl

    Get specific acl with 'content'='config':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=sample-ipv4-acl?content=config

    Get specific acl with 'content'='nonconfig':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=sample-ipv4-acl?content=nonconfig

    Get specific acl with 'content'='all':
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=sample-ipv4-acl?content=all

Remove Filtering Rules

    $ ./do_request_from_file.sh DELETE {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=sample-ipv4-acl

## DB

To set up your database, refer to the [Database configuration document](./docs/DATABASE.md)  
The 'dots_server' accesses the 'dots' database on MySQL as the root user.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ mysql -u root -p dots < ./dots_server/db_models/test_dump.sql

Or you can run MySQL on docker.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ docker run -d -p 3306:3306 -v ${PWD}/dots_server/db_models/test_dump.sql:/docker-entrypoint-initdb.d/test_dump.sql:ro -e MYSQL_DATABASE=dots -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql

DOTS server listens to DB notification (e.g changes to mitigation_scope#status) at port 9999. If you want to change to different port, you have to change it at two places:

    - dots_server/dot_server.yaml: dbNotificationPort: 9999
    - mysql_udf/mysql-notification.c: #define PORT 9999

After changing port number, it is neccessary to rebuild go-dots (which includes rebuilding mysql-notification.c and restarting DB) so that the change can take effect.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ make && make install

# GOBGP

Check the route is installed successfully in gobgp server

    $ gobgp global rib

```
Network              Next Hop             AS_PATH              Age        Attrs
*> 172.16.238.100/32    172.16.238.254                            00:00:42   [{Origin: i}]
```
