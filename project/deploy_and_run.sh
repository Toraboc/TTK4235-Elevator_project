#!/usr/bin/env bash

set -euo pipefail

# ---------------- CONFIG ----------------
IP_PREFIX="10.100.23."
REMOTE_USER="student"
LOCAL_GO_DIR="."
REMOTE_BASE_DIR="/home/$REMOTE_USER/Documents/Sanntid55/"
CODE_DIR="."
SSH_KEY="$HOME/.ssh/id_ed25519"
SSH_PUB_KEY="$SSH_KEY.pub"
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

if [[ "$(uname -s)" == "Linux" && ! -f "$SSH_KEY" ]]; then
    echo "No SSH key found. Generating one..."
    ssh-keygen -t ed25519 -f "$SSH_KEY" -N ""
fi

ensure_ssh_access() {
    local host="$1"

    echo "Checking SSH access to $host..."

    if ssh -o BatchMode=yes -o ConnectTimeout=5 "$REMOTE_USER@$host" "echo ok" &>/dev/null; then
        echo "Key-based SSH already works for $host"
        return
    fi

    echo "Key-based SSH not available for $host"
    echo "You will be prompted for the SSH password."

    ssh-copy-id -i "$SSH_PUB_KEY" "$REMOTE_USER@$host"

    echo "SSH key installed on $host"
}


# Keep track of background job PIDs
declare -a JOB_PIDS=()

# Ctrl-C handler
cleanup() {
    echo "Caught Ctrl-C! Killing background jobs..."

    # Kill local background jobs (SSH sessions)
    for pid in "${JOB_PIDS[@]}"; do
        disown "$pid" 2>/dev/null || true
        kill "$pid" 2>/dev/null || true
    done

    # Kill go run . and elevatorserver on all remote hosts
    for last_octet in "${HOSTS[@]}"; do
        host="${IP_PREFIX}${last_octet}"
        echo "Stopping remote processes on $host"
        ssh "$REMOTE_USER@$host" "pkill -f 'go run .'" || true
        ssh "$REMOTE_USER@$host" "pkill -f 'project'" || true
    done

    exit 1
}

trap cleanup INT

run_remote() {
    local host="$1"

    echo "Deploying to $host..."

    # Kill running elevatorserver, to make sure all other controllers stop
    ssh "$REMOTE_USER@$host" "
        pkill -f ./elevatorserver
    " || true

    ssh "$REMOTE_USER@$host" "mkdir -p $REMOTE_BASE_DIR/$CODE_DIR"
    rsync -a --delete "$CODE_DIR" "$REMOTE_USER@$host:$REMOTE_BASE_DIR/$CODE_DIR/"

    ssh "$REMOTE_USER@$host" "
        set -e

        cd $REMOTE_BASE_DIR$CODE_DIR
        echo 'Starting elevatorserver'
        ./elevatorserver > elevatorserver.log 2>&1 &
        sleep 1

        go run . 2>&1
    " | sed "s/^/[$host] /" &
    JOB_PIDS+=("$!")
}

for last_octet in "${HOSTS[@]}"; do
    [[ -z "$last_octet" ]] && continue

    host="${IP_PREFIX}${last_octet}"

    ensure_ssh_access "$host"
    run_remote "$host"
done

wait
