#!/bin/bash
set -ex
openssl genrsa -aes256 -passout pass:1234 -out grafarg-e2e-ca.key.pem 4096
openssl rsa -passin pass:1234 -in grafarg-e2e-ca.key.pem -out grafarg-e2e-ca.key.pem
openssl req -config openssl.cnf -key grafarg-e2e-ca.key.pem -new -x509 -days 3650 -sha256 -extensions v3_ca -out grafarg-e2e-ca.pem