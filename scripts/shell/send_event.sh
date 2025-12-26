#!/bin/bash

EVENT_ID=$(uuidgen)
SOURCE="sh-script"
PAYLOAD='{"msg":"shell test"}'

curl -s -o /dev/null -w "%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"event_id\":\"$EVENT_ID\",\"source\":\"$SOURCE\",\"payload\":$PAYLOAD}" \
  http://localhost:8080/v1/events
