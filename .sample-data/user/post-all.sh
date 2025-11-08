#!/usr/bin/env sh

for userData in data/*.json; do
  curl -X POST localhost:3333/api/auth/signup \
    -H "Content-Type: application/json" \
    --data-binary "@$userData"
done
