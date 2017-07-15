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



# Usage

## Server
    $ $GOPATH/bin/dots_server -host 127.0.0.1 -port 4646 -file [path to the config file]

## Client
    $ $GOPATH/bin/dots_client -host 127.0.0.1 -port 4646 
    
## Client Controller [mitigation_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Post \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json
   
## Client Controller [mitigation_cancel_request]
    $ $GOPATH/bin/dots_client_controller -request mitigation_request -method Delete \
     -json $GOPATH/src/github.com/nttdots/go-dots/dots_client/sampleMitigationRequest.json

# Example

## CentOS7 on Docker [mitigation_request]
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



