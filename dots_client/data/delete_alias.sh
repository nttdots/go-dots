#!/bin/sh

HOST_PATH=$1
CUID=$2
NAME=$3

P="`dirname $0`"

if [ -z "$HOST_PATH" ]; then
  echo 'No HOST_PATH parameter.' >&2
  exit 1
fi
if [ -z "$CUID" ]; then
  echo 'No CUID specified.' >&2
  exit 1
fi
if [ -z "$NAME" ]; then
  echo 'No NAME specified.' >&2
  exit 1
fi

$P/do_request.sh \
  DELETE \
  $HOST_PATH/data/ietf-dots-data-channel:dots-data/dots-client="$CUID"/aliases/alias="$NAME"
