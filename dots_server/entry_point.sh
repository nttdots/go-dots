#!/usr/bin/env bash

/usr/bin/wait-for-it.sh -t 30 db:3306 || exit 1

(sleep 3 && nc -l 8080) &

$@
