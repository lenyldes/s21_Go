#!/bin/bash
set -e

cert_dir="$(cd -- "$(dirname -- "$0")" && pwd)/certs"
mkdir -p "$cert_dir"

# Самоподписанный сертификат для локального домена be3.shop на один год.
# Новая генерация заменяет существующие сертификат и ключ.
openssl req -x509 -newkey rsa:2048 -sha256 -nodes \
    -keyout "$cert_dir/be3.shop.key" \
    -out "$cert_dir/be3.shop.crt" \
    -days 365 \
    -subj "/CN=be3.shop" \
    -addext "subjectAltName=DNS:be3.shop"

chmod 600 "$cert_dir/be3.shop.key"
