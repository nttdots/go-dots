#!/bin/bash

/etc/quagga/init_conf.sh

ipv4=`ip a show eth0|grep inet |awk '{print $2}' | sed 's/\/.*$//'| head -n1`

sed -i "s/ROUTER_IP_V4/${ipv4}/" /root/gobgpd.conf

gobgpd -f /root/gobgpd.conf &

#for debugging
bash


