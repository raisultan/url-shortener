#!/bin/bash

set -e

host="$1"
port="$2"
dep_name="$3"
shift 3
cmd="$@"

until PING="$(timeout 1 nc -z -v $host $port 2>&1)"; do
  echo "$dep_name is unavailable - sleeping"
  sleep 3
done

echo "Dependencies are up - executing command"
exec $cmd
