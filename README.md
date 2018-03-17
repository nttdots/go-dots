

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
* make, autoconf, automake, libtool
* [git](https://git-scm.com/)
* [go](https://golang.org/doc/install)
  * go 1.9 or later is required. (for latest GoBGP)
  * set PATH to go and set $GOPATH, using their instructions.
* [openssl](https://www.openssl.org/)
  * OpenSSL 1.1.0 or higher (for libcoap)

* MySQL 5.7 or higher

## Recommandation Environment
* Ubuntu 16.04+
* macOS High Sierra 10.13+

## How to build go-dots
### build libcoap custom for go-dots
    
    $ wget https://github.com/nttdots/go-dots/blob/ietf101interop/misc/libcoap_custom_for_go-dots.tar.gz
    $ tar zxvf libcoap_custom_for_go-dots.tar.gz
    $ cd libcoap_custom_for_go-dots
    $ ./autogen.sh
    $ ./configure --disable-documentation --with-openssl=/usr/local
    $ make
    $ sudo make install
    
### Install go-dots
To install go-dots source codes and command line programs, use the following command:
    
    $ go get -u github.com/nttdots/go-dots/...
    

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
    $ $GOPATH/bin/dots_server -config [config.yml file (ex: go-dots/dots_server/dots_server.yaml)]


## Client
    $ $GOPATH/bin/dots_client --server localhost --signalChannelPort=5684 -vv

    
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

## DB

To set up your database, refer to the [Database configuration document](./docs/DATABASE.md)  
The 'dots_server' accesses the 'dots' database on MySQL as the root user.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ mysql -u root dots < ./dots_server/db_models/test_dump.sql

Or you can run MySQL on docker.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ docker run -d -p 3306:3306 -v ${PWD}/dots_server/db_models/test_dump.sql:/docker-entrypoint-initdb.d/test_dump.sql:ro -e MYSQL_DATABASE=dots -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql


# GOBGP

Check the route is installed successfully in gobgp server

    $ docker exec -it gobgp gobgp global rib

```
Network              Next Hop             AS_PATH              Age        Attrs
*> 172.16.238.100/32    172.16.238.254                            00:00:42   [{Origin: i}]
```
