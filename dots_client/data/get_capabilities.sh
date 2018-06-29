#!/bin/sh

HOST_PATH=$1
CONTENT=$2

P="`dirname $0`"

if [ -z "$HOST_PATH" ]; then
  echo 'No HOST_PATH parameter.' >&2
  exit 1
fi
if [ -n "$CONTENT" ]; then
  CONTENT_PARAM="content=$CONTENT"
fi

$P/do_request.sh \
  GET \
  $HOST_PATH/data/ietf-dots-data-channel:dots-data/capabilities?"$CONTENT_PARAM"
