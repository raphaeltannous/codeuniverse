#!/usr/bin/env sh

curl -X POST "localhost:3333$1" \
  -H "Content-Type: application/json" \
  --data-binary "@$2"
