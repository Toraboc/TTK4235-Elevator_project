#!/usr/bin/env bash

set -euo pipefail

# ---------------- CONFIG ----------------
IP_PREFIX="10.100.23."
REMOTE_USER="student"
LOCAL_GO_DIR="."
REMOTE_BASE_DIR="/home/$REMOTE_USER/Documents/Sanntid55"
CODE_DIR="/project"
SSH_KEY="$HOME/.ssh/id_ed25519"
SSH_PUB_KEY="$SSH_KEY.pub"
# ----------------------------------------

# Parse flags
FAST_MODE=false
HOSTS=()

while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -f)
            FAST_MODE=true
            shift
            ;;
        *)
            HOSTS+=("$1")
            shift
            ;;
    esac
done

if [[ "${#HOSTS[@]}" -eq 0 ]]; then
    echo "Usage: $0 [-f] host1 [host2 ...]"
    echo "  -f: Fast mode (skip SSH setup and elevatorserver)"
    echo "For example: ./deploy_and_run.sh -f 33 34 35"
    exit 1
fi

# 1. Ensure SSH key exists
if [[ "$(uname -s)" == "Linux" && ! -f "$SSH_KEY" ]]; then
    echo "No SSH key found. Generating one..."
    ssh-keygen -t ed25519 -f "$SSH_KEY" -N ""
fi

# 2. Function to ensure key-based access
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
    
    # Kill local background jobs
    for pid in "${JOB_PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
    done

    # Kill elevatorserver on all remote hosts
    for last_octet in "$@"; do
        host="${IP_PREFIX}${last_octet}"
        echo "Stopping elevatorserver on $host"
        ssh "$REMOTE_USER@$host" "pkill -f './elevatorserver'" || true
    done

    exit 1
}

trap cleanup INT

# 3. Copy code and start program (fast mode)
run_remote_fast() {
    local host="$1"

    echo "Fast deploying to $host..."

    # Copy code
    scp -r ".$CODE_DIR" "$REMOTE_USER@$host:$REMOTE_BASE_DIR"

    # Start go code only
    ssh "$REMOTE_USER@$host" "
        set -e
        cd $REMOTE_BASE_DIR$CODE_DIR
        go run . 2>&1
    " | sed "s/^/[$host] /" &
    JOB_PIDS+=("$!")
}

# 4. Copy code and start program (full mode)
run_remote() {
    local host="$1"

    echo "Deploying to $host..."

    # 1. Stop elevatorserver if something is listening on port 15657
    ssh "$REMOTE_USER@$host" "
        pkill -f ./elevatorserver
    " || true

    # Copy code
    ssh "$REMOTE_USER@$host" "
        rm -r $REMOTE_BASE_DIR
        mkdir -p $REMOTE_BASE_DIR
    " || true

    scp -r ".$CODE_DIR" "$REMOTE_USER@$host:$REMOTE_BASE_DIR"
    scp elevatorserver "$REMOTE_USER@$host:$REMOTE_BASE_DIR"

    # Start elevatorserver and go code
    ssh "$REMOTE_USER@$host" "
        set -e

        cd $REMOTE_BASE_DIR
        echo 'Starting elevatorserver'
        ./elevatorserver > elevatorserver.log 2>&1 &
        sleep 1

        cd $REMOTE_BASE_DIR$CODE_DIR
        go run . 2>&1
    " | sed "s/^/[$host] /" &
    JOB_PIDS+=("$!")
}

# 5. Main loop
for last_octet in "${HOSTS[@]}"; do
    [[ -z "$last_octet" ]] && continue

    # Construct full IP
    host="${IP_PREFIX}${last_octet}"

    if [[ "$FAST_MODE" == false ]]; then
        ensure_ssh_access "$host"
        run_remote "$host"
    else
        run_remote_fast "$host"
    fi
done

wait
