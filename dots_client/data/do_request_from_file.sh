#!/bin/sh

CERTS_DIR=$1
METHOD=$2
URI_PATH=$3
FILE=$4

if [ -z "$METHOD" ]; then
  echo 'No METHOD parameter.' >&2
  exit 1
fi
if [ -z "$URI_PATH" ]; then
  echo 'No URI_PATH parameter.' >&2
  exit 1
fi
if { [ -n "$FILE" ] && ! [ -e "$FILE" ]; }; then
  echo "Input file $FILE is not existed"
  exit 1
fi

if [ -n "$FILE" ]; then
wget \
  -q -S -O - \
  --content-on-error \
  --no-check-certificate \
  --ca-certificate="$CERTS_DIR"/ca-cert.pem \
  --certificate="$CERTS_DIR"/client-cert.pem \
  --private-key="$CERTS_DIR"/client-key.pem \
  --method="$METHOD" \
  "$URI_PATH" \
  "--body-file=$FILE" \
  "--header=Content-Type: application/yang-data+json"
else
wget \
  -q -S -O - \
  --content-on-error \
  --no-check-certificate \
  --ca-certificate="$CERTS_DIR"/ca-cert.pem \
  --certificate="$CERTS_DIR"/client-cert.pem \
  --private-key="$CERTS_DIR"/client-key.pem \
  --method="$METHOD" \
  "$URI_PATH"
fi
