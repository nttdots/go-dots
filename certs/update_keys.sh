#!/usr/bin/env bash

CERTTOOL=certtool

# update CA
$CERTTOOL --generate-privkey --bits 4096 --outfile ca-key.pem
$CERTTOOL --generate-self-signed --load-privkey ca-key.pem --template template_ca.txt --outfile ca-cert.pem
$CERTTOOL --generate-crl --load-ca-privkey ca-key.pem --load-ca-certificate ca-cert.pem --template template_ca.txt --outfile crl.pem
$CERTTOOL --generate-dh-params --outfile dh.pem --sec-param medium

# update server
$CERTTOOL --generate-privkey --bits 4096 --outfile server-key.pem
$CERTTOOL --generate-request --load-privkey server-key.pem --template template_server.txt --outfile server-csr.pem
$CERTTOOL --generate-certificate --load-request server-csr.pem --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem  --template template_server.txt --outfile server-cert.pem
rm -f server-csr.pem

# update client
$CERTTOOL --generate-privkey --bits 4096 --outfile client-key.pem
$CERTTOOL --generate-request --load-privkey client-key.pem --template template_client.txt --outfile client-csr.pem
$CERTTOOL --generate-certificate --load-request client-csr.pem --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem  --template template_client.txt --outfile client-cert.pem
rm -f client-csr.pem
