#!/bin/sh

CUID=$1

P="`dirname $0`"

if [ -z "$CUID" ]; then
  echo 'No CUID specified.' >&2
  exit 1
fi

$P/do_request.sh \
  GET \
  restconf/data/ietf-dots-data-channel:dots-data/dots-client="$CUID"/aliases
