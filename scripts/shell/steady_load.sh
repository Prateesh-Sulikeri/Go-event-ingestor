#!/bin/bash

RPS=${1:-5}

echo "Sending $RPS requests per second... CTRL+C to stop."

while true; do
  for ((i=1; i<=$RPS; i++)); do
    ./send_event.sh &
  done
  wait
  sleep 1
done
