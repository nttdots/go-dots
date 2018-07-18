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
case "$NAME" in
  "https1")
    RULE='{ "ietf-dots-data-channel:aliases": { "alias": [ { "name": "https1", "target-protocol": [ 6 ], "target-prefix": [ "2001:db8:6401::1/128", "2001:db8:6401::2/128" ], "target-port-range": [ { "lower-port": 443 } ] } ] } }'
    ;;
  "Server1")
    RULE='{ "ietf-dots-data-channel:aliases": { "alias": [ { "name": "Server1", "target-protocol": [ 6 ], "target-prefix": [ "2001:db8:6401::1/128", "2001:db8:6401::2/128" ], "target-port-range": [ { "lower-port": 443 } ] } ] } }'
    ;;
  "Server2")
    RULE='{ "ietf-dots-data-channel:aliases": { "alias": [ { "name": "Server2", "target-protocol": [ 6 ], "target-prefix": [ "2001:db8:6401::10/128", "2001:db8:6401::20/128" ], "target-port-range": [ { "lower-port": 80 } ] } ] } }'
    ;;
  *)
    echo 'Unknown alias name, use "https1", "Server1" or "Server2".' >&2
    exit 1
    ;;
esac

$P/do_request.sh \
  PUT \
  $HOST_PATH/data/ietf-dots-data-channel:dots-data/dots-client="$CUID"/aliases/alias="$NAME" \
  "$RULE"

