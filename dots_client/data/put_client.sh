#!/bin/sh

HOST_PATH=$1
CUID=$2
CDID=$3

P="`dirname $0`"

if [ -z "$HOST_PATH" ]; then
  echo 'No HOST_PATH parameter.' >&2
  exit 1
fi
if [ -z "$CUID" ]; then
  echo 'No CUID specified.' >&2
  exit 1
fi
if [ -n "$CDID" ]; then
  $P/do_request.sh \
    PUT \
    $HOST_PATH/data/ietf-dots-data-channel:dots-data/dots-client="$CUID" \
    '{ "ietf-dots-data-channel:dots-client": [ {"cuid": "'"$CUID"'", "cdid": "'"$CDID"'" } ] }'
else
  $P/do_request.sh \
    PUT \
    $HOST_PATH/data/ietf-dots-data-channel:dots-data/dots-client="$CUID" \
    '{ "ietf-dots-data-channel:dots-client": [ {"cuid": "'"$CUID"'" } ] }'
fi

