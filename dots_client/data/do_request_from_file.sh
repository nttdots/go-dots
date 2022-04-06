#!/bin/bash
CA_CERT=../certs/ca-cert.pem
CLIENT_CERT=../certs/client-cert.pem
CLIENT_KEY=../certs/client-key.pem
while [[ $# -gt 0 ]]; do
    case $1 in
        --ca-cert)
            CA_CERT=$2
            shift 2
            ;;
        --client-cert)
            CLIENT_CERT=$2
            shift 2
            ;;
        --client-key)
            CLIENT_KEY=$2
            shift 2
            ;;
        POST|PUT)
            METHOD=$1
            URI_PATH=$2
            FILE=$3
            shift 3
            ;;
        GET|DELETE)
            METHOD=$1
            URI_PATH=$2
            shift 2
            ;;
        *)
            echo "Unknown option $1"
            exit 1
            ;;
    esac
done

if { [ -n "$CA_CERT" ] && ! [ -e "$CA_CERT" ]; }; then
    echo "$CA_CERT is not existed"
    exit 1
fi
if { [ -n "$CLIENT_CERT" ] && ! [ -e "$CLIENT_CERT" ]; }; then
    echo "$CLIENT_CERT is not existed"
    exit 1
fi
if { [ -n "$CLIENT_KEY" ] && ! [ -e "$CLIENT_KEY" ]; }; then
    echo "$CLIENT_KEY is not existed"
    exit 1
fi

if { [ -n "$FILE" ] && ! [ -e "$FILE" ]; }; then
    echo "$FILE is not existed"
    exit 1
fi

if [ -n "$FILE" ]; then
wget \
  -q -S -O - \
  --content-on-error \
  --no-check-certificate \
  --ca-certificate=$CA_CERT \
  --certificate=$CLIENT_CERT \
  --private-key=$CLIENT_KEY \
  --method="$METHOD" \
  "$URI_PATH" \
  "--body-file=$FILE" \
  "--header=Content-Type: application/yang-data+json"
else
wget \
  -q -S -O - \
  --content-on-error \
  --no-check-certificate \
  --ca-certificate=$CA_CERT \
  --certificate=$CLIENT_CERT \
  --private-key=$CLIENT_KEY \
  --method="$METHOD" \
  "$URI_PATH"
fi