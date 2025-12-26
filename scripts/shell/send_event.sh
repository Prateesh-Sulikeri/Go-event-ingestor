#!/bin/bash
DIR=$(dirname "$0")

TOKEN=$(bash "$DIR/get_token.sh")

EVENT_ID=$(uuidgen)

curl -s -o /dev/null -w "%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"event_id\":\"$EVENT_ID\",\"source\":\"dashboard\",\"payload\":{\"origin\":\"ui\"}}" \
  http://localhost:8080/v1/events
