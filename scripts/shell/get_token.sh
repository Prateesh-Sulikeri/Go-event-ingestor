#!/bin/bash

CLIENT_ID="${1:-dashboard}"

TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"client_id\":\"$CLIENT_ID\"}" | jq -r '.token')

if [[ "$TOKEN" == "null" || -z "$TOKEN" ]]; then
  echo "Error: Could not generate token"
  exit 1
fi

echo "$TOKEN"
