#!/usr/bin/env bash

set -euo pipefail

# ---------------- CONFIG ----------------
MACHINES_FILE="machines.txt"
REMOTE_USER="student"
LOCAL_GO_DIR="$HOME/Documents/Sanntid15"
REMOTE_BASE_DIR=$LOCAL_GO_DIR
GO_MAIN="main.go"
SSH_KEY="$HOME/.ssh/id_ed25519"
SSH_PUB_KEY="$SSH_KEY.pub"
# ----------------------------------------

if [[ ! -f "$MACHINES_FILE" ]]; then
    echo "Machines file not found: $MACHINES_FILE"
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

    ssh "$REMOTE_USER@$host" "mkdir -p ~/.ssh && chmod 700 ~/.ssh"
    scp "$SSH_PUB_KEY" "$REMOTE_USER@$host:/tmp/authorized_keys.tmp"

    ssh "$REMOTE_USER@$host" "
        cat /tmp/authorized_keys.tmp >> ~/.ssh/authorized_keys &&
        chmod 600 ~/.ssh/authorized_keys &&
        rm /tmp/authorized_keys.tmp
    "

    echo "SSH key installed on $host"
}

# 3. Copy code and start program
run_remote() {
    local host="$1"

    echo "Deploying to $host..."

    ssh "$REMOTE_USER@$host" "mkdir -p $REMOTE_BASE_DIR"
    scp -r "$LOCAL_GO_DIR"/* "$REMOTE_USER@$host:$REMOTE_BASE_DIR/"

    ssh "$REMOTE_USER@$host" "
        cd $REMOTE_BASE_DIR &&
        go run $GO_MAIN 2>&1
    " | sed "s/^/[$host] /"
}

# 4. Main loop
while read -r host; do
    [[ -z "$host" ]] && continue

    ensure_ssh_access "$host"

    # Run each host in background so outputs interleave
    run_remote "$host" &
done < "$MACHINES_FILE"

wait
echo "All deployments started."
