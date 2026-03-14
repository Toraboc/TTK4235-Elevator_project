# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
cd project && go build .

# Run locally (builds + supervises both processes)
./supervisor.sh

# Run directly
cd project && go build . && ./project -server localhost:15657

# Docker (3 simulators + 3 controllers)
docker compose up --build

# Deploy to lab machines
./deploy_and_run.sh 35 36 37   # last octets of 10.100.23.x
```

No test suite exists. No linter is configured.

## Architecture

Three independent goroutines communicate via channels. Each process owns its state exclusively — no shared memory between them.

```
Elevator Process  ←→  Order Handler Process  ←→  Network Process
```

**Elevator** (`elevator/`) — hardware interface. Polls buttons/floor sensors, drives motor, manages door timer. Publishes state changes and completed orders. Receives target floors and confirmed orders (for lights).

**Order Handler** (`orderHandler/`) — distributed coordination. Maintains `WorldView`: all elevators' states + orders across the network. On each event, recomputes which orders to assign using an external D-language binary (`hall_request_assigner`) and publishes new target floors.

**Network** (`network/`) — UDP broadcast at 100 Hz to `255.255.255.255:44043`. Each message carries sender ID, elevator state, orders, and known peers. Nodes are pruned after 200 ms of silence. Only counts a node as connected if both sides know each other.

## Key Data Structures

**WorldView** — owned by order handler, no mutex needed:
- `Orders map[NodeId]*Orders` — hall up/down + cab orders per node
- `ConnectedNodes NodeIdSet` — active nodes right now
- `ElevatorStates map[NodeId]ElevatorState` — floor, direction, behaviour per node

**Order status** is a cyclic state machine: `NO_ORDER → UNCONFIRMED → CONFIRMED → FINISHED → NO_ORDER`. Advancement requires agreement across connected nodes, giving eventual consistency during network splits.

**NodeId** is a `uint32` encoding of the IPv4 address.

## Important Conventions

- Channel interfaces use directional types (`chan<-` / `<-chan`). Request/response uses a reply channel sent over a channel.
- `panic()` is used for unrecoverable state (motor fault timeout, impossible elevator state). This is intentional — the supervisor restarts the process.
- `hall_request_assigner` binary (D language, precompiled) must be in the `project/` working directory. It is called via stdin/stdout JSON.
- `elevatorserver` binary (precompiled, in repo root) is the hardware simulator, listening on port 15657 by default.
- The `shared/` package holds types used across all three processes — avoid adding process-specific logic there.
