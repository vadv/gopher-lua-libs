#!/usr/bin/env bash
#
# Make a test cert and key

readonly MY_DIR="$(dirname "${BASH_SOURCE[0]}")"

openssl req -nodes -x509 -newkey rsa:4096 -keyout "${MY_DIR}/test.key.pem" -out "${MY_DIR}/test.cert.pem" -sha256 \
  -days 365000 -addext 'subjectAltName = IP:127.0.0.1'
