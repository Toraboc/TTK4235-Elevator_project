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

Use `./project -h` for all flags.

## Architecture

Three concurrent processes communicate via channels, each owning their state exclusively:

- **Elevator** — hardware interface, motor/door control state machine
- **Order Handler** — maintains world view across all nodes, assigns orders via `hall_request_assigner`
- **Network** — UDP broadcast for peer discovery and state syn
