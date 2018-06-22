#!/bin/sh

METHOD=$1
URI_PATH=$2
BODY=$3

if [ -z "$METHOD" ]; then
  echo 'No METHOD parameter.' >&2
  exit 1
fi
if [ -z "$URI_PATH" ]; then
  echo 'No URI_PATH parameter.' >&2
  exit 1
fi

CERTS_DIR="`dirname $0`/../../certs"

if [ -n "$BODY" ]; then
wget \
  -q -S -O - \
  --no-check-certificate \
  --ca-certificate="$CERTS_DIR"/ca-cert.pem \
  --certificate="$CERTS_DIR"/client-cert.pem \
  --private-key="$CERTS_DIR"/client-key.pem \
  --method="$METHOD" \
  "https://127.0.0.1:10443/$URI_PATH" \
  "--body-data=$BODY" \
  "--header=Content-Type: application/yang-data+json"
else
wget \
  -q -S -O - \
  --no-check-certificate \
  --ca-certificate="$CERTS_DIR"/ca-cert.pem \
  --certificate="$CERTS_DIR"/client-cert.pem \
  --private-key="$CERTS_DIR"/client-key.pem \
  --method="$METHOD" \
  "https://127.0.0.1:10443/$URI_PATH"
fi
