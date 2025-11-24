#!/usr/bin/env sh

curl -X POST localhost:3333/api/auth/signup \
  -H "Content-Type: application/json" \
  --data-binary "@$1"
