#!/bin/sh

# Make root certificate

openssl ecparam -out root.key -name prime256v1 -genkey
openssl req -new -sha256 -key root.key -out root.csr -subj "/C=RU/ST=Moscow/O=Yandex-Practicum/CN=root-monogo"
openssl x509 -req -sha256 -days 365 -in root.csr -signkey root.key -out root.crt
openssl x509 -in root.crt -out root.pem -outform PEM

# Make server certificate

openssl ecparam -out server.key -name prime256v1 -genkey
openssl req -new -sha256 -key server.key -out server.csr -subj "/C=RU/ST=Moscow/O=Yandex-Practicum/CN=server-monogo"
openssl x509 -req -in server.csr -CA root.crt -CAkey root.key -CAcreateserial -out server.crt -days 365 -sha256

# Make agent certificate

openssl ecparam -out agent.key -name prime256v1 -genkey
openssl req -new -sha256 -key agent.key -out agent.csr -subj "/C=RU/ST=Moscow/O=Yandex-Practicum/CN=agent-monogo"
openssl x509 -req -in agent.csr -CA root.crt -CAkey root.key -CAcreateserial -out agent.crt -days 365 -sha256
