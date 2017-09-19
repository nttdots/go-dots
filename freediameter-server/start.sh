#!/bin/sh

/usr/bin/wait-for-it.sh -t 30 diameap_db:3306 || exit 1

/usr/bin/freeDiameterd -d -c /freeDiameter.conf
