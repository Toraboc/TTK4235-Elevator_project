# Elevator Controller

Distributed elevator controller for TTK4145. Coordinates multiple elevators over a local network with fault tolerance and automatic order reassignment.

## Requirements

- Go 1.24+
- `elevatorserver` binary in this directory
- `hall_request_assigner` binary in this directory

## Run

```bash
# Managed (recommended) — builds and auto-restarts on crash
./supervisor.sh

# Direct
# Remember to start the elevatorserver before the go program.
go run .
```

## Architecture

Three concurrent processes communicate via channels, each owning their state exclusively:

- **Elevator** — hardware interface, motor/door control state machine
- **Order Handler** — maintains world view across all nodes, assigns orders via `hall_request_assigner`
- **Network** — UDP broadcast for peer discovery and state sync

## Deploy script

Use `deploy_and_run.sh` to build and deploy to lab machines.
Machines are identified by the last octet of their IP address.
Below is an example of usage.

```bash
./deploy_and_run.sh 35 36 37
```

**First deployment to a new machine:** deploy to it alone, not in parallel with others.
The first connection copies SSH keys and requires a password, which is difficult to enter while other deployments are running simultaneously.
