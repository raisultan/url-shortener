#!/bin/bash

set -e

host="$1"
shift
cmd="$@"

until PING="$(timeout 1 nc -z -v $host 9000 2>&1)"; do
  echo "ClickHouse is unavailable - sleeping"
  sleep 1
done

echo "Dependencies are up - executing command"
exec $cmd
