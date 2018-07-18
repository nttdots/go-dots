#!/bin/sh

HOST_PATH=$1

if [ -z "$HOST_PATH" ]; then
  echo 'No HOST_PATH parameter.' >&2
  exit 1
fi
P="`dirname $0`"
CERTS_DIR="`dirname $0`/../../certs"

$P/do_request.sh \
  GET \
  $HOST_PATH/.well-known/host-meta

