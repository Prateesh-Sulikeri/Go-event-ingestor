#!/bin/bash
COUNT=${1:-200}

echo "Sending burst of $COUNT requests..."

for i in $(seq 1 $COUNT); do
  ./send_event.sh &
done

wait
echo "Burst complete."
