#!/bin/sh

echo "Enter Country (2 letter code):"
read COUNTRY
echo "Enter State or Province Name (full name):"
read STATE
echo "Enter Locality Name (eg, city):"
read LOCALITY
echo "Enter Organization Name (eg, company):"
read ORGANIZATION
echo "Enter Organizational Unit Name (eg, section):"
read ORG_UNIT
echo "Enter Common Name (e.g. server FQDN or YOUR name):"
read COMMON_NAME

openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
  -keyout ./certs/server.key -out ./certs/server.cert \
  -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=$ORG_UNIT/CN=$COMMON_NAME"