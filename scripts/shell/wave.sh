#!/bin/bash
DIR=$(dirname "$0")
MIN=${1:-1}
MAX=${2:-20}

trap "echo 'Stopping wave'; pkill -P $$; exit 0" SIGINT

while true; do
  for i in $(seq "$MIN" "$MAX"); do
    echo "Wave: $i rps"
    bash "$DIR/steady_load.sh" "$i" &
    sleep 1
    pkill -P $$ steady_load.sh
  done
  for i in $(seq "$MAX" -1 "$MIN"); do
    echo "Wave: $i rps"
    bash "$DIR/steady_load.sh" "$i" &
    sleep 1
    pkill -P $$ steady_load.sh
  done
done
