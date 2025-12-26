#!/bin/bash

MAX=${1:-20}
MIN=${2:-1}

while true; do
  for i in $(seq $MIN $MAX); do
    echo "RPS $i"
    ./steady_load.sh $i &
    sleep 1
    pkill -f steady_load.sh
  done

  for i in $(seq $MAX -1 $MIN); do
    echo "RPS $i"
    ./steady_load.sh $i &
    sleep 1
    pkill -f steady_load.sh
  done
done
