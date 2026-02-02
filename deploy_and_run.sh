#!/usr/bin/env bash

set -euo pipefail

# ---------------- CONFIG ----------------
IP_PREFIX="10.100.23."
REMOTE_USER="student"
LOCAL_GO_DIR="$HOME/Documents/Sanntid55"
REMOTE_BASE_DIR=$LOCAL_GO_DIR
CODE_DIR="/project"
GO_MAIN="main.go"
SSH_KEY="$HOME/.ssh/id_ed25519"
SSH_PUB_KEY="$SSH_KEY.pub"
# ----------------------------------------

if [[ "$#" -eq 0 ]]; then
    echo "Usage: $0 host1 [host2 ...]"
    echo "For example ./deploy_and_run.sh 33 34 35"
    exit 1
fi

# 1. Ensure SSH key exists
if [[ ! -f "$SSH_KEY" ]]; then
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

# 3. Copy code and start program
run_remote() {
    local host="$1"

    echo "Deploying to $host..."

    # 1. Stop elevatorserver if something is listening on port 15657
    ssh "$REMOTE_USER@$host" "
        pkill -f ./elevatorserver
    " || true

    # Copy code
    ssh "$REMOTE_USER@$host" "mkdir -p $REMOTE_BASE_DIR"
    scp -r "$LOCAL_GO_DIR"/* "$REMOTE_USER@$host:$REMOTE_BASE_DIR/"

    # Start elevatorserver and go code
    ssh "$REMOTE_USER@$host" "
        set -e

        cd $REMOTE_BASE_DIR
        echo 'Starting elevatorserver'
        ./elevatorserver > elevatorserver.log 2>&1 &
        sleep 1
        
        cd $REMOTE_BASE_DIR$CODE_DIR
        go run $GO_MAIN 2>&1
    " | sed "s/^/[$host] /"
}

# 4. Main loop
for last_octet in "$@"; do
    [[ -z "$last_octet" ]] && continue

    # Construct full IP
    host="${IP_PREFIX}${last_octet}"

    ensure_ssh_access "$host"

    # Run each host in background so outputs interleave
    run_remote "$host" &
done


wait
