fdd8
7
- The main file is readable and clearly shows how the different modules are set up and the interfaces between them.
- The modules are clear, and most of the implementations are where they are expected. However, there are some coherence issues, for example, cab lights and hall lights are controlled from different modules. Storing cab orders to a file as a backup is also a responsibility shared between the orderHandler and elevatorcontrol modules.
- All modules have a nice, minimal public interface, but several implementation details are not properly hidden. For example, other modules must know that elevatorcontrol needs some time to start up, this leaks an internal detail. See the sleep call in the orderHandler module.
- The state machine in elevatorcontrol is readable, and it looks easy to add functionality.
- Some variable names could be more consistent. For example, constants in the elevatorcontrol module could better describe what they represent, a "rate" implies a frequency, not a duration, and it is not obvious that "buffer" refers to a time duration.
- Global variables are used in some modules where they do not seem necessary, and they introduce side effects that are visible across the entire module.
- A lot of the communication between modules and goroutines involves polling channels for updates, when it would be simpler to just block and wait for data to arrive. This also applies to elevatorcontrol, where a global variable is polled for state changes instead of waiting on a channel. This would also reduce CPU usage.
- There are several places where the code could be simplified and dead code removed (code that runs but has no effect). For example, in processPair, a goroutine is started and the function immediately waits for it to finish, which is unnecessary. Additionally, many functions end with an if/else and an early return, the else branch is redundant when the if block already returns.
- The use of goto can make code harder to read. The instances where it appears could be simplified into a function with an early return.
- In the network module, raw syscalls are used to set up the UDP socket, which is lower level than necessary. A package would make it easier to read.

dce1
9
- All modules are imported with a dot import, which makes it harder to immediately see which module a function belongs to when reading the code.
- The main file clearly shows the three modules and how they are connected. Channels are created here and passed through typed interface structs, so both the components and their dependencies are visible at a glance. The one oddity is the busy-waiting for loop at the end, the last goroutine could just be run as a regular function call instead.
- From the file structure alone you can find where each design decision lives: hallRequestAssigner.go for order assignment, cyclicCounter.go for the distributed confirmation protocol, targetFloor.go for how the elevator decides where to go next, and nodeControl.go for peer discovery.
- Each process owns its state exclusively and it is always clear who is responsible for what. Shared state across goroutines is essentially absent, with myId being the only global variable, and that is write-once at startup.
- In the orderHandler, the worldView methods are the only functions that perform side effects, which makes it easy to reason about where state changes actually happen.
- In the cyclic counter, the fieldSelector parameter could have a more descriptive name. Passing a function to select a field in a struct can be a sign of poor struct organisation, but here the structs are designed to be convenient when working elevator by elevator rather than button by button, so a selector function is actually a useful solution to a real structural tension.
- Names like WorldView, ConfirmedOrders, transferPassengers, and cyclicCounter communicate intent clearly and never mislead you when navigating the codebase.

e012
7
- Code is nicely split into modules, but there are a few too many, making the codebase feel a bit chaotic. The "wvm" module should be given a more descriptive name ("world view manager" is not obvious from the acronym).
- The code does not build, there are unused imports and variables that break compilation ("fmt" in fsm/fsm.go, request/request.go, wvm/wvm.go). Some code also looks commented out last-minute, which makes the snapshot hard to run/review.
- There are many comments that should likely be removed, especially self-notes and comments describing unused code. Keeping these adds noise and makes the code harder to scan. Prefer self-documenting code and keep comments for non-obvious decisions. Also remember to remove comments about removing comments(request.go) 🙂
- Network broadcast unnecessarily uses a low-level abstraction with OS-specific code. Go's net package can handle UDP broadcast cross-platform, which would likely simplify the code and improve readability.
- Most functions are declared global even when only used internally. Consider unexporting helper functions (lowercase) to reduce coupling and make it clearer which functions are intended to be used by other modules.
- Overall structure is coherent, but responsibility boundaries are sometimes unclear (who owns/updates which state). Tightening ownership per module would improve readability and testability.
- Enums and structs are centralized in types.go, but several would be easier to navigate if placed closer to their owning module (or at least grouped more consistently); there are also inconsistencies in how this is handled across the codebase.
- Function naming could be clearer. For example, network transmitters (peers.Transmitter vs bcast.Transmitter) don't communicate intent at call sites; more specific names and/or a small shared abstraction would make the code easier to follow.

bf32
6
- You have a event-driven structure with clear goroutines for driver, network, and scheduler, but the structure is fairly problematic. All routines are run from elevatorManager.go instead of main.go. In addition, for instance the fsm module gives light, motor and timer side effects, and shows clear signs of weak cohesion. The elevatormodule does both all the major operations and supersmall ones. To fix the structure it would be easier to restart than to refactor.
- ElevatorManager works well as orchestrator, but consider renaming or moving scheduler-specific helpers (for example publishSchedulerSnapshot) to keep module responsibilities cleaner. In addition the whole managerState is passed to several modules, which only should need small parts of the state, and this weakens cohesion.
- Request handling already has useful structure, but clarifying ownership of Requests, Assignments, and GlobalCabRequests would make future changes safer.
- Better naming and a cleaner file stucture could make the code more readable; for instance the name RequestStates is a bit unclear, and the map showing the states is fairly hidden in the filestructure.
- Network and registry interaction is functional, but defining a sharper boundary for peer liveness vs. transport responsibilities would simplify reasoning.
- It might be benefitial to add a short architecture section in README describing acknowledgement, order assignment, execution flow, and backup behavior, so design choices are visible without deep code tracing.

