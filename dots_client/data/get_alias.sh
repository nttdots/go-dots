#!/bin/sh

CUID=$1
NAME=$2
CONTENT=$3

P="`dirname $0`"

if [ -z "$CUID" ]; then
  echo 'No CUID specified.' >&2
  exit 1
fi
if [ -z "$NAME" ]; then
  echo 'No NAME specified.' >&2
  exit 1
fi

if [ -n "$CONTENT" ]; then
  CONTENT_PARAM="content=$CONTENT"
fi

$P/do_request.sh \
  GET \
  restconf/data/ietf-dots-data-channel:dots-data/dots-client="$CUID"/aliases/alias="$NAME"?"$CONTENT_PARAM"
