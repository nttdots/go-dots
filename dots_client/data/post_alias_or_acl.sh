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
  "sample-ipv4-acl")
    RULE='{ "ietf-dots-data-channel:acls": { "acl": [ { "name": "sample-ipv4-acl", "type": "ipv4-acl-type", "activation-type": "activate-when-mitigating", "aces": { "ace": [ { "name": "rule1", "matches": { "ipv4": { "destination-ipv4-network": "198.51.100.0/24", "source-ipv4-network": "192.0.2.0/24" } }, "actions": { "forwarding": "drop" } } ] } } ] } }'
     ;;
  "dns-fragments-ipv4")
    RULE='{ "ietf-dots-data-channel:acls": { "acl": [ { "name": "dns-fragments-ipv4", "type": "ipv4-acl-type", "aces": { "ace": [ { "name": "drop-all-except-last-fragment", "matches": { "ipv4": { "flags": "more" } }, "actions": { "forwarding": "drop" } }, { "name": "allow-dns-packets", "matches": { "ipv4": { "destination-ipv4-network": "198.51.100.0/24" }, "udp": { "destination-port": { "operator": "eq", "port": 53 } } }, "actions": { "forwarding": "accept" } }, { "name": "drop-last-fragment", "matches": { "ipv4": { "flags": "more" } }, "actions": { "forwarding": "drop" } } ] } } ] } }'
    ;;
  "dns-fragments-ipv6")
    RULE='{ "ietf-dots-data-channel:acls": { "acl": [ { "name": "dns-fragments-ipv6", "type": "ipv6-acl-type", "aces": { "ace": [ { "name": "drop-all-fragments", "matches": { "ipv6": { "fragment": [null] } }, "actions": { "forwarding": "drop" } }, { "name": "allow-dns-packets", "matches": { "ipv6": { "destination-ipv6-network": "2001:db8::/32" }, "udp": { "destination-port": { "operator": "eq", "port": 53 } } }, "actions": { "forwarding": "accept" } } ] } } ] } }'
    ;;
  *)
    echo 'Unknown name, use "https1", "Server1", "Server2", "sample-ipv4-acl", "dns-fragments-ipv4" or "dns-fragments-ipv6".' >&2
    exit 1
    ;;
esac

$P/do_request.sh \
  POST \
  $HOST_PATH/data/ietf-dots-data-channel:dots-data/dots-client="$CUID" \
  "$RULE"

