#!/bin/sh

CONTENT=$1

P="`dirname $0`"

if [ -n "$CONTENT" ]; then
  CONTENT_PARAM="content=$CONTENT"
fi

$P/do_request.sh \
  GET \
  restconf/data/ietf-dots-data-channel:dots-data/capabilities?"$CONTENT_PARAM"
