# Install

## Required Software

* gcc 4.8.5 or later is required.
* make
* [git](https://git-scm.com/)
* [go](https://golang.org/doc/install)
  * go 1.6 or later is required.
  * set PATH to go and set $GOPATH, using their instructions.
* [gnutls](http://www.gnutls.org/)
  * gnutls 3.3.24 or later is required.

* MySQL 5.7 or later is required.

### Issue 

* go 1.8 build failed with macOS/Xcode 8.3.

    go 1.8.1 will resolve the problem.
    
    another solution is to install Xcode8.2
    
    [runtime: some Go executables broken with macOS 10.12.4 / Xcode 8.3 but fine with Xcode 8.2 ](https://github.com/golang/go/issues/19734)
    
## go-dots

To install go-dots source and command line program, use the following:

    $ go get -u github.com/nttdots/go-dots/...


# How it works

go-dots consists of 3 daemons below:

* dots_server: dots server agent
* dots_client: dots client specified in the Internet drafts
* dots_client_controller: The controller for the dots_client. Operators controll the dots_client by this controller

To explain the relationships between each daemons in detail, here we illustrates the mitigation request procedure of our system. First, dots_client_controller passes request JSONs to the dots_client. Examples of these JSONs are located in the 'dots_client/' directory in this Github project. Remember to edit the 'mitigation-id' fields in the JSON file every time you make requests to the dots_server. Then the dots client converts the received JSON to the CBOR format and send a mitigation request to the server via a CoAP connection.

The figure below shows the detailed sequence diagram which depicts how a dots_server handle a mitigation request. 

<img src='https://github.com/nttdots/go-dots/blob/documentation/docs/pics/mitigation_request_sequence.png' title='mitigation requests sequence diagram'>

Upon receiving CoAP messages, a dots_server first find the appropriate message controller for the received message. The message controller first find the customer information bound to the common name field contained in the client certificate. The customer information is configured and stored in the RDB before the server is started. If no appropriate customer objets is found, the dots_server decline the received request. 

Then the server retrieves mitigation scopes contained in the received mitigation request message. Again the dots_server validate the mitigation scopes to determine whether the customer which issued the request has the valid privilege for the mitigation operations. If the validation is successfully completed, the server select a blocker for the mitigation and execute the mitigation.

# Usage

## Server
    $ $GOPATH/bin/dots_server -config [path to the config.yml file(ex: go-dots/dots_server/dots_server.yaml)]

## Client
    $ $GOPATH/bin/dots_client -host 127.0.0.1 -port 4646 
    
## Client Controller [mitigation_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Post \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json
   
## Client Controller [mitigation_cancel_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Delete \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json

# One box example on Docker

## mitigation_request

### build and up gobgp server
    $ cd $GOPATH/src/github.com/nttdots/go-dots/gobgp-server/
    $ docker-compose build
    $ docker-compose up -d

### build and up dots client/server and db
    $ cd $GOPATH/src/github.com/nttdots/go-dots/example/mitigation-request
    $ docker-compose build
    $ docker-compose up

# Test

## Prepare

### SetUp DB

dots accesses dots_test as root on MySQL.

You should import dump.sql for test data.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ mysql -u root dots_test < ./dots_server/db_models/test_dump.sql


Or you can run MySQL on docker.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ docker run -d -p 3306:3306 -v ${PWD}/dots_server/db_models/test_dump.sql:/docker-entrypoint-initdb.d/test_dump.sql:ro -e MYSQL_DATABASE=dots_test -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql

### Run tests

You can test all package.

    $ cd $GOPATH/src/github.com/nttdots/go-dots/
    $ make test



