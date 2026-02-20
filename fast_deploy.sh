#!/usr/bin/env bash

set -euo pipefail

# ---------------- CONFIG ----------------
IP_PREFIX="10.100.23."
REMOTE_USER="student"
REMOTE_BASE_DIR="/tmp"
CODE_DIR="project"
# ----------------------------------------
HOSTS=()

while [[ "$#" -gt 0 ]]; do
    HOSTS+=("$1")
    shift
done

if [[ "${#HOSTS[@]}" -eq 0 ]]; then
    echo "Usage: $0 host1 [host2 ...]"
    echo "For example: ./deploy_and_run.sh 33 34 35"
    exit 1
fi

cd $CODE_DIR
GOOS=linux GOARCH=amd64 go build .
cd ..
mv "$CODE_DIR/project" elevator

run_remote() {
    local host="$1"

    scp elevator "$REMOTE_USER@$host:$REMOTE_BASE_DIR"

    ssh "$REMOTE_USER@$host" "$REMOTE_BASE_DIR/elevator" | sed "s/^/[$host] /"
}

# # Ctrl-C handler
# cleanup() {
#     echo "Caught Ctrl-C! Killing background jobs..."
#
#     exit 1
# }
#
# trap cleanup INT

for last_octet in "${HOSTS[@]}"; do
    [[ -z "$last_octet" ]] && continue

    # Construct full IP
    host="${IP_PREFIX}${last_octet}"
    run_remote "$host" &
done

wait
