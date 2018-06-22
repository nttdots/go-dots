#!/bin/sh

CUID=$1
CDID=$2

P="`dirname $0`"

if [ -z "$CUID" ]; then
  echo 'No CUID specified.' >&2
  exit 1
fi
if [ -n "$CDID" ]; then
  $P/do_request.sh \
    PUT \
    restconf/data/ietf-dots-data-channel:dots-data/dots-client="$CUID" \
    '{ "ietf-dots-data-channel:dots-client": [ {"cuid": "'"$CUID"'", "cdid": "'"$CDID"'" } ] }'
else
  $P/do_request.sh \
    PUT \
    restconf/data/ietf-dots-data-channel:dots-data/dots-client="$CUID" \
    '{ "ietf-dots-data-channel:dots-client": [ {"cuid": "'"$CUID"'" } ] }'
fi

