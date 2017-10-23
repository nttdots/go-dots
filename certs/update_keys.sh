#!/usr/bin/env bash

CERTTOOL=gnutls-certtool

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

# update invalid client
$CERTTOOL --generate-privkey --bits 4096 --outfile invalid_client-key.pem
$CERTTOOL --generate-request --load-privkey invalid_client-key.pem --template template_invalid_client.txt --outfile invalid_client-csr.pem
$CERTTOOL --generate-certificate --load-request invalid_client-csr.pem --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem  --template template_invalid_client.txt --outfile invalid_client-cert.pem
rm -f invalid_client-csr.pem

# update not_auth client
$CERTTOOL --generate-privkey --bits 4096 --outfile not_auth_client-key.pem
$CERTTOOL --generate-request --load-privkey not_auth_client-key.pem --template template_not_auth_client.txt --outfile not_auth_client-csr.pem
$CERTTOOL --generate-certificate --load-request not_auth_client-csr.pem --load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem  --template template_not_auth_client.txt --outfile not_auth_client-cert.pem
rm -f not_auth_client-csr.pem
