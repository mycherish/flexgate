#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
THREADS="${THREADS:-4}"
CONNECTIONS="${CONNECTIONS:-100}"
DURATION="${DURATION:-30s}"

echo "Benchmark users route"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency \
  "$BASE_URL/api/users/1"

echo "Benchmark orders route"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency \
  "$BASE_URL/api/orders"