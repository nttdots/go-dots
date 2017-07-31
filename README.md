

# go-dots

![logo](https://github.com/nttdots/go-dots/blob/master/go-dots_logo/go-dots_logo_blue.png)

This is a DDoS Open Threat Signaling (dots) implementation written in Go. This implmentation is based on the Internet drafts below. 

* draft-ietf-dots-architecture-04 
* draft-ietf-dots-data-channel-02 
* draft-ietf-dots-requirements-06 
* draft-ietf-dots-signal-channel-02 
* draft-ietf-dots-use-cases-06 

This implementation is not fully compliant with the documents listed above.  For example, we are utilizing CoAP as the data channel protocol while the current version of the data channel document specifies RESTCONF as the data channel protocol.

Licensed under Apache License 2.0.

# How to Install

## Requirements

* gcc 4.8.5 or higher
* make
* [git](https://git-scm.com/)
* [go](https://golang.org/doc/install)
  * go 1.6 or later is required.
  * set PATH to go and set $GOPATH, using their instructions.
* [gnutls](http://www.gnutls.org/)
  * gnutls 3.3.24 or higher

* MySQL 5.7 or higher

### An issue to build this project on MacOS

* Building go 1.8 fails with macOS/Xcode 8.3.

    go 1.8.1 will resolve the problem.
    
    Or you can solve the problem by installing Xcode8.2
    
    [runtime: some Go executables broken with macOS 10.12.4 / Xcode 8.3 but fine with Xcode 8.2 ](https://github.com/golang/go/issues/19734)
    
## How to build go-dots

To install go-dots source codes and command line programs, use the following command:

    $ go get -u github.com/nttdots/go-dots/...


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

Or,

    $ cd $GOPATH/src/github.com/nttdots/go-dots/example/dots_server
    $ docker-compose build
    $ docker-compose up 

## Client
    $ $GOPATH/bin/dots_client -server 127.0.0.1

Or,

    $ cd $GOPATH/src/github.com/nttdots/go-dots/example/dots_client
    $ docker-compose build
    $ docker-compose up
    
### Client Controller [mitigation_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Post \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json
   
### Client Controller [mitigation_cancel_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Delete \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json

## DB

To set up your database, refer to the [Database configuration document](./docs/DATABASE.md)  
The 'dots_server' accesses the 'dots' database on MySQL as the root user.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ mysql -u root dots < ./dots_server/db_models/test_dump.sql

Or you can run MySQL on docker.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ docker run -d -p 3306:3306 -v ${PWD}/dots_server/db_models/test_dump.sql:/docker-entrypoint-initdb.d/test_dump.sql:ro -e MYSQL_DATABASE=dots -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql


# Test

## One box example on Docker (mitigation request and delete)

Build dots client, server, db and gobgp in one box and connect them each other on a docker network.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/example/onebox
    $ docker-compose build
    $ docker-compose up

Setup customer information in db

    $ docker exec -i db mysql -u root dots < customer_example.sql

You can see how they work by this example command on the dots_client.

    $ docker exec -i dots_client dots_client_controller -method Post -request mitigation_request -json example/onebox/sampleMitigationRequest.json

Check the route is installed successfully in gobgp server

    $ docker exec -it gobgp gobgp global rib

```
Network              Next Hop             AS_PATH              Age        Attrs
*> 172.16.238.100/32    172.16.238.254                            00:00:42   [{Origin: i}]
```
You can withdraw the route by Delete method

    $ docker exec -i dots_client dots_client_controller -method Delete -request mitigation_request -json example/onebox/sampleMitigationRequest.json

You can restore the db

    $ docker exec -i db mysql -u root dots < ../../dots_server/db_models/template.sql 


## Running tests

You can test all package by the commands below.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ make test



