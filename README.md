

# go-dots

![logo](https://github.com/nttdots/go-dots/blob/master/go-dots_logo/go-dots_logo_blue.png)

"go-dots" is a DDoS Open Threat Signaling (dots) implementation written in Go. This implmentation is based on the Internet drafts below. 

* RFC9132 (was draft-ietf-dots-rfc8782-bis-08)
* RFC 8783 (was draft-ietf-dots-data-channel)
* draft-ietf-dots-architecture-18
* RFC 8612 (was draft-ietf-dots-requirements)
* draft-ietf-dots-use-cases-21
* draft-ietf-dots-signal-filter-control-07
* draft-ietf-dots-signal-call-home-09
* draft-ietf-dots-telemetry-16
* draft-ietf-dots-robust-blocks-00

This implementation is not fully compliant with the documents listed above.  For example, we are utilizing CoAP as the data channel protocol while the current version of the data channel document specifies RESTCONF as the data channel protocol.

Licensed under Apache License 2.0.

# How to Install

## Requirements

* gcc 5.4.0 or higher
* make, autoconf, automake, libtool, pkg-config, pkgconf or pkg-config
* [git](https://git-scm.com/)
* [go](https://golang.org/doc/install)
  * go 1.13.5 or later is required. (for the latest GoBGP - v2.20.0)
  * set PATH to go and set $GOPATH, using their instructions.
* [openssl](https://www.openssl.org/)
  * OpenSSL 1.1.1g or higher (for libcoap)

* MySQL 5.7.x and its development package (MySQL 8.0.x or higher not yet supported)
  * Install mysql development package in Ubuntu:
    $ sudo apt-get install libmysqld-dev

* gnuTLS (install to configure the certificate)
    $ sudo apt-get install gnutls-bin

## Recommandation Environment
* Ubuntu 16.04+
* macOS High Sierra 10.13+

## How to build go-dots
### Build libcoap for go-dots

To build libcoap for go-dots. We will work as follow:

- Pull libcoap. Currenly supported libcoap version : v4.3.0
    ```
    $ git clone https://github.com/obgm/libcoap.git
    $ cd libcoap
    $ git checkout a80d462ff57630ce214efdf5caf34133b02ad7ee

- Merge [q-block](https://github.com/mrdeep1/libcoap/tree/q-block) into libcoap.
    ```
    $ cd libcoap
    $ sudo git remote add remote_name https://github.com/mrdeep1/libcoap.git
    $ sudo git pull remote_name q-block

- Build libcoap.
    ```
    $ cd libcoap
    $ ./autogen.sh
    $ ./configure --disable-documentation --with-openssl
    $ sudo make && sudo make install

### Install gorilla-mux for go-dots client router to handle RESTful API

    $ go get -u github.com/gorilla/mux
    
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

The blocker configuration of DOTS server is defined in database. For more detail about database, refer to the [Database configuration document](./docs/DATABASE.md)

## Server
    $ $GOPATH/bin/dots_server --config [config.yml file (ex: go-dots/dots_server/dots_server.yaml)]


## Client
    $ $GOPATH/bin/dots_client --server localhost --signalChannelPort=4646 --config [config.yml file (ex: go-dots/dots_client/dots_client.yaml)] -vv

## MySQL Notification
After the go-dots is built, the mysql-notification.* MUST be copied from $GOPATH/src/github.com/nttdots/go-dots/mysql-udf to /usr/lib/mysql/plugin

    $ sudo cp -avr $GOPATH/src/github.com/nttdots/go-dots/mysql-udf/* /usr/lib/mysql/plugin

## GoBGP Server
To install and run gobgp-server, refer to the following link:
* [gobgp-server](https://github.com/osrg/gobgp)


## Arista Server
The Arista Server (Arista Switch Hardware) is installed and executed on Arista Extensible Operating System (EOS) platform.
In order to connect to Arista Server, DOTS server use 'goeapi' open source library that is provided by Arista Networks company and the configuration file (.eapi.conf) which contains information about the Arista networks such as: host, username, password, etc.

```
[connection:arista]
    host=arista
    username=admin
    password=123456
    enablepwd=passwd
    transport=https
```

For more detailed information about the configuration of 'goeapi', refer to the following link:
* [arista-goeapi](https://github.com/aristanetworks/goeapi)


##  Signal Channel
The primary purpose of the signal channel is for a DOTS client to ask a DOTS server for help in mitigating an attack, and for the DOTS server to inform the DOTS client about the status of such mitigation.

### Client Controller [mitigation_request]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraft.json

In order to handle out-of-order delivery of mitigation requests, 'mid' values MUST increase monotonically. Besides, if the 'mid' value has exceeded 3/4 of (2**32 - 1), it should be reset by sending a mitigation request with 'mid' is set to '0' to avoid 'mid' rollover. However, the reset request is only accepted by DOTS server at peace-time (have no any active mitigation request which is maintaining).

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=0 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraft.json

### Client Controller [mitigation_retrieve_all]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw

### Client Controller [mitigation_retrieve_one]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123

### Client Controller [mitigation_retrieve_all_query]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -targetProtocol=17 -aliasName=https1

### Client Controller [mitigation_retrieve_one_query]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -targetPrefix=1.2.0.10/32

### Client Controller [mitigation_withdraw]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Delete \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123

### Client Controller [mitigation_observe]
A DOTS client can convey the 'observe' option set to '0' in the GET request to receive notification whenever status of mitigation request changed and unsubscribe itself by issuing GET request with 'observe' option set to '1'

Subscribe for resource observation:

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -observe=0

Unsubscribe from resource observation:

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Get \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -observe=1

Subscriptions are valid as long as current session exists. When session is renewed (e.g DOTS client does not receive response from DOTS server for its Ping message in a period of time, it decided that server has been disconnected, then re-connects), all previous subscriptions will be lost. In such cases, DOTS clients will have to re-subscribe for observation. Below is recommended step: 

    ・GET a list of all existing mitigations (that were created before server restarted)
    ・PUT mitigations  one by one
    ・GET + Observe for all the mitigations that should be observed

### Client Controller [mitigation_efficacy_update]
A DOTS client can convey the 'If-Match' option with empty value in the PUT request to transmit DOTS mitigation efficacy update to the DOTS server:

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -ifMatch="" \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraftEfficacyUpdate.json

DOTS client to DOTS server mitigation efficacy DOTS telemetry attributes

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 -ifMatch="" \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationEfficacyTelemetryAttributes.json

### Client Controller [session_configuration_request]

    $ $GOPATH/bin/dots_client_controller -request session_configuration -method Put \
     -sid 234 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleSessionConfigurationDraft.json

In order to handle out-of-order delivery of session configuration, 'sid' values MUST increase monotonically.

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
The primary purpose of the data channel is to support DOTS related configuration and policy information exchange between the DOTS client and the DOTS server.

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
    $ ./do_request_from_file.sh POST {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleAlias.json

    Put alias:
    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/aliases/alias=xxx sampleAlias.json

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

    Post acl with insert:
    $ ./do_request_from_file.sh POST '{href}/data/ietf-dots-data-channel:dots-data/dots-client=123?insert=first' sampleAcl.json

    Post acl with insert and point:
    $ ./do_request_from_file.sh POST '{href}/data/ietf-dots-data-channel:dots-data/dots-client=123?insert=after&point=xxx' sampleAcl.json

    Put acl:
    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=xxx sampleAcl.json

    Put acl with insert:
    $ ./do_request_from_file.sh PUT '{href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=xxx?insert=last' sampleAcl.json

    Put acl with insert and point:
    $ ./do_request_from_file.sh PUT '{href}/data/ietf-dots-data-channel:dots-data/dots-client=123/acls/acl=xxx?insert=before?point=xxx1' sampleAcl.json
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

### Managing Vendor Attack Mapping
Create vendor-mapping

    $ ./do_request_from_file.sh POST {href}/data/ietf-dots-data-channel:dots-data/dots-client=123 sampleVendorAttackMapping.json

Update vendor-mapping

    $ ./do_request_from_file.sh PUT {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping/vendor-id=345 sampleVendorAttackMapping.json

Get vendor-mapping of server

    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/ietf-dots-mapping:vendor-mapping

Get vendor-mapping

    Get vendor-mapping with 'depth'
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping?depth=3

    Get vendor-mapping with 'content'
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping?content=all

    Get vendor-mapping without 'depth' and 'content'
    $ ./do_request_from_file.sh GET {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping?content=all

Delete vendor-mapping

    Delete one vendor-mapping
    $ ./do_request_from_file.sh DELETE {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping/vendor-id=345

    Delete all vendor-mapping
    $ ./do_request_from_file.sh DELETE {href}/data/ietf-dots-data-channel:dots-data/dots-client=123/ietf-dots-mapping:vendor-mapping

## Signal Channel Control Filtering
Unlike the DOTS signal channel, the DOTS data channel is not expected to deal with attack conditions.
Therefore, when DOTS client is under attacked by DDoS, the DOTS client can use DOTS signal channel protocol to manage the filtering rule in DOTS Data Channel to enhance the protection capability of DOTS protocols.

### Client Controller [mitigation_control_filtering]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraftControlFiltering.json

## Signal Channel Call Home
The DOTS signal channel Call Home identify the source to block DDoS attack traffic closer to the source(s) of a DDoS attack.
when the DOTS client is under attacked by DDoS, the DOTS client sends the attack traffic information to the DOTS server. The DOTS server in turn uses the attack traffic information to identify the compromised devices launching the outgoing DDoS attack and takes appropriate mitigation action.

### Client Controller [mitigation_call_home]

    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Put \
     -cuid=dz6pHjaADkaFTbjr0JGBpw -mid=123 \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequestDraftCallHome.json

##  Telemetry
The telemetry aims to enrich DOTS signal channel protocol with various telemetry attributes allowing optimal DDoS attack mitigation. The telemetry specifies the normal traffic baseline and attack traffic telemetry attributes a DOTS client can convey to its DOTS server in the mitigation request, the mitigation status telemetry attributes a DOTS server can communicate to a DOTS client, and the mitigation efficacy telemetry attributes a DOTS client can communicate to a DOTS server. The telemetry contains `Telemetry Setup Configuration` and `Telemetry Pre-or-ongoing-mitigation`


### Client Controller [telemmetry_setup_request] (Telemetry Setup Configuration)
Registering telemetry configuration

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Put -cuid=dz6pHjaADkaFTbjr0JGBpw -tsid=123\
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleTelemetryConfiguration.json

Registering total pipe capacity

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Put -cuid=dz6pHjaADkaFTbjr0JGBpw -tsid=123\
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleTotalPipeCapacity.json

Registering baseline

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Put -cuid=dz6pHjaADkaFTbjr0JGBpw -tsid=123\
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleBaseline.json

Get one telemetry setup configuration

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw -tsid=123

Get all telemetry setup configuration

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw

Delete one telemetry setup configuration

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Delete -cuid=dz6pHjaADkaFTbjr0JGBpw -tsid=123

Delete all telemetry setup configuration

    $ $GOPATH/bin/dots_client_controller -request telemetry_setup_request -method Delete -cuid=dz6pHjaADkaFTbjr0JGBpw

### Client Controller [telemmetry_pre_mitigation_request] (Telemetry Pre-Or-Ongoing-Mitigation)
Registering telemetry pre-or-ongoing-mitigation

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Put -cuid=dz6pHjaADkaFTbjr0JGBpw -tmid=123\
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleTelemetryPreMitigation.json

Get one telemetry pre-or-ongoing-mitigation

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw -tmid=123

Get all telemetry pre-or-ongoing-mitigation

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw

Get one telemetry pre-or-ongoing-mitigation with query

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw -tmid=123 -targetProtocol=17
    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw -tmid=123 -c=a

Get all telemetry pre-or-ongoing-mitigation with query

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Get -cuid=dz6pHjaADkaFTbjr0JGBpw -targetProtocol=17

Delete one telemetry pre-or-ongoing-mitigation

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Delete -cuid=dz6pHjaADkaFTbjr0JGBpw -tmid=123

Delete all telemetry pre-or-ongoing-mitigation

    $ $GOPATH/bin/dots_client_controller -request telemetry_pre_mitigation_request -method Delete -cuid=dz6pHjaADkaFTbjr0JGBpw

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

# Arista
Check the active Arista ACL that is installed successfully in Arista server

    $ show running-config interfaces Ethernet 1

```
interface Ethernet1
   description DDoS-Seg-to_internet
   ip address 172.16.238.100/30
   ip access-group mitigation-acl-1 in
```

Check the active rules in Arista ACL that is installed successfully in Arista server

    $ show ip access-lists mitigation-acl-1

```
IP Access List mitigation-acl-1
        10 deny 4 any 1.1.2.0/24
        20 deny 4 any host 1.1.1.69
```

# GOBGP Flowspec

Check the flowspec route is installed successfully in gobgp server

    $ gobgp global rib -a ipv4-flowspec
    $ gobgp global rib -a ipv6-flowspec

```
   Network                                                                   Next Hop      AS_PATH    Age        Attrs
*> [destination: 1.1.2.0/24][protocol: ==tcp][destination-port: >=443&<=800] fictitious               00:00:06   [{Origin: i} {Extcomms: [redirect: 1.1.1.0:100]}]
```

# Certificate Configuration

### Precondition
* The GnuTLS has been installed.

### Configure The Certificate

Typically, the client's/server's certificate is a single identifier type which means that the certificate has only one common name (CN-ID) as identifier. However, the Common Name is not strongly typed because the Common Name can contain a human-friendly string, not a DNS domain name. Moreover, the client's/server's certificate can be multi identifiers type which including the Common Name (CN-ID) and some Subject Alternative Name (DNS-ID, SRV-ID), in order to ensure that at least one DNS qualified domain name. In go-dots, two certificate types are configured as below:

* The single identifier type

   The template file only has the Common Name (CN-ID)

   Example: In file [template_client](./certs/template_client.txt), configure as follow:
    ```
    # X.509 server certificate options

    organization = "sample client"
    state = "Tokyo"
    country = JP
    cn = "client.sample.example.com"
    expiration_days = 365

    # X.509 v3 extensions
    signing_key
    encryption_key
    ```

* The multi identifiers type

    The template file has the Common Name (CN-ID) and some Subject Alternative Name (DNS-ID, SRV-ID)

    Example: In file [template_client](./certs/template_client.txt), add more Subject Alternative Name as follow:

    ```
    # DNS name(s) of the server
    dns_name = "xample1.example.com"
    dns_name = "_xampp.example.com
    ```

To add/change the server's certificate or the client's certificate, execute with following command:
```
./update_keys.sh
```

For more detailed information about configuration of the certificate and the GnuTLS, refer to the following link:
* [action_gnutls_scripted.md](https://gist.github.com/epcim/832cec2482a255e3f392)

# Kubernetes

Kubernetes is a portable, extensible, open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation. It has a large, rapidly growing ecosystem. Kubernetes services, support, and tools are widely available.

**To install kubernetes, refer to the following link:**
* [kubernetes-ubuntu](https://matthewpalmer.net/kubernetes-app-developer/articles/install-kubernetes-ubuntu-tutorial.html)

**To setup and run the go-dots server and go-dots client on Kubernetes. Following as below:**
- Starting the virtual machine minikube

    ```
    $ minikube start
    ```

- Created the dots folder to Mount it to minikube. The dots folder contains the certs folder and the config folder. The certs folder contains the cert files in the [certs](./certs) folder. The config folder contains the [test_dump.sql](./dots_server/db_models/test_dump.sql), the gobgpd.conf, the [dots_server.yaml](./dots_server/dots_server.yaml) and the [dots_client.yaml](./dots_client/dots_client.yaml).
    The structure of the dots folder as below:
    ```
    ~/dots
        |--/certs
                |-- /certs/* files
        |--/config
                |-- /test_dump.sql
                |-- /gobgpd.conf
                |-- /dots_server.yaml
                |-- /dots_client.yaml
    ```

- Mounting the dots folder from host directory to minikube. Then the host directory and the virtual machine minikube contain the dots folder with the files and the values of files that are same.
    ```
    $ minikube mount ~/dots:/dots
    ```

- Created the go-dots client, the go-dots server, the mysql and the gobgp on Kuberbetes by the [Pod.yaml](./docker/Pod.yaml) file
    ```
    $ kubectl create -f Pod.yaml
    ```

- Get Pod to check ip address of the go-dots server and the go-dots client
    ```
    $ kubectl get pods --output=wide
    ```

- Executing the container in the cluster by specifying the Pod name
    - The go-dots client
        ```
        $ kubectl exec -it client /bin/bash
        ```
    - The go-dots server
        ```
        $ kubectl exec -it server /bin/bash
        ```
    - The mysql
        ```
        $ kubectl exec -it server -c mysql /bin/bash
        ```
    - The gobgp
        ```
        $ kubectl exec -it server -c gobgp /bin/bash
        ```

- After executing the go-dots server container. In the go-dots server container, you to do as below: 
    - Building the go-dots
        ```
        $ cd $GOPATH/src/github.com/nttdots/go-dots
        $ make && make install
        ```
    - Copy the mysql_udf/* to /usr/lib/mysql/plugin to the go-dots server listens to DB notification
        ```
        $ sudo cp $GOPATH/src/github.com/nttdots/go-dots/mysql_udf/* /usr/lib/mysql/plugin
        ```
    - Run the go-dots server
        ```
        $ cd /dots/config
        $ $GOPATH/bin/dots_server --config dots_server.yaml -vv
        ```

- After executing the mysql container. In the mysql container, restored the `dots` database from the mysql dump file
    ```
    $ cd /dots/config
    $ mysql -u root -proot dots < test_dump.sql
    ```

- After executing the gobgp container. In the gobgp container, run gobgp:
    ```
    $ cd /dots/config
    $ sudo $GOPATH/bin/gobgpd -f gobgpd.conf
    ```

- After executing the go-dots client container. In the go-dots client container, run the go-dots client:
    ```
    $ cd /dots/config
    $ $GOPATH/bin/dots_client --server `ip address of the go-dots server` --signalChannelPort=4646 --config dots_client.yaml -vv
    ```

- If you want to change port (signalchannel, datachannel) to different port in dots_server.yaml, you have to change it at two places:
    - ~/dots/config/dot_server.yaml: signalChannelPort: 4646
    - ~/dots/config/dot_server.yaml: dataChannelPort: 10443

- If you want to change port (signalchannel, datachannel) to different port in Pod.yaml, you have to change it at two places:
    - [/docker/Pod.yaml](./docker/Pod.yaml): containerPort: 4646
    - [/docker/Pod.yaml](./docker/Pod.yaml): containerPort: 10443

- If you want to change the fields in dots_server.yaml, you have to change it at '~/dots/config/dots_server.yaml'

- If you want to insert new dots_client, you to do as below:
    - In the mysql container, Insert new dots_client into database
        ```
        $ mysql -u root -proot
        >mysql use dots;
        >mysql INSERT INTO `customer` (`id`, `common_name`, `created`, `updated`) VALUES (129,'client.sample.com','2017-04-13 13:44:34','2017-04-13 13:44:34'),
        ```
    - Next, In the '~/dots/certs/' folder add client certification as the `Certificate Configuration` above