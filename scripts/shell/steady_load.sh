#!/bin/bash
DIR=$(dirname "$0")
RPS=${1:-5}

echo "Steady Load: $RPS requests per second"
while true; do
  for _ in $(seq 1 "$RPS"); do
    bash "$DIR/send_event.sh" &
  done
  wait
  sleep 1
done
