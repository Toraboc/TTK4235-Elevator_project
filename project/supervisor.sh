#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ELEVATORSERVER="$SCRIPT_DIR/elevatorserver"
BINARY="$SCRIPT_DIR/project"

cd "$SCRIPT_DIR" && go build -o "$BINARY" .

killJobs() {
    pkill -f "elevatorserver" 2>/dev/null || true
}

cleanup() {
    echo "Shutting down..."
    killJobs
    exit 0
}

trap cleanup INT TERM

while true; do
    killJobs
    sleep 0.5

    echo "Starting elevatorserver..."
    "$ELEVATORSERVER" &
    ELEVATOR_PID=$!
    sleep 1

    echo "Starting application..."
    "$BINARY" &
    APP_PID=$!

    wait -n "$ELEVATOR_PID" "$APP_PID" 2>/dev/null || true

    echo "A process exited, restarting both..."
done
