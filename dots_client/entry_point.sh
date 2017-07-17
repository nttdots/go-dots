#!/bin/bash

# /usr/bin/wait-for-it.sh -t 30 ${DOTS_SERVER_IPV4}:8080 || exit 1

dots_client -vv -server ${DOTS_SERVER_IPV4} -signalChannelPort ${SIGNAL_CHANNEL_PORT} -dataChannelPort ${DATA_CHANNEL_PORT} -clientCertFile ${GOPATH}/src/github.com/nttdots/go-dots/certs/client-cert.pem -certFile ${GOPATH}/src/github.com/nttdots/go-dots/certs/ca-cert.pem -clientKeyFile ${GOPATH}/src/github.com/nttdots/go-dots/certs/client-key.pem
while :
do
    sleep 1
done

$@
