#!/bin/bash

CLIENT_ID=${1:-"tester"}

TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"client_id\":\"$CLIENT_ID\"}" | jq -r '.token')

echo "Exporting TOKEN..."
export TOKEN
echo $TOKEN
