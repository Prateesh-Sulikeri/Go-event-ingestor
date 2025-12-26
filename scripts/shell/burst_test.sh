#!/bin/bash
DIR=$(dirname "$0")
COUNT=${1:-200}

echo "Burst: $COUNT events"
for _ in $(seq 1 "$COUNT"); do
  bash "$DIR/send_event.sh" &
done
wait
