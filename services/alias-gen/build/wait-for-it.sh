#!/bin/bash

 set -e

 host="$1"
 shift
 cmd="$@"

 until PING="$(timeout 1 nc -z -v $host 5432 2>&1)"; do
   echo "PostgreSQL is unavailable - sleeping"
   sleep 1
 done

 echo "Dependencies are up - executing command"
 exec $cmd
