import core.sync.semaphore;
import core.time;

// Initial values:
// M     = 1
// PS[2] = [0, 0]
// busy  = false
// numWaiting = [0, 0]

__gshared Semaphore M;
__gshared Semaphore[2] PS;
__gshared bool busy = false;
__gshared int[2] numWaiting = [0, 0];

shared static this() {
	M = new Semaphore(1);
	PS[0] = new Semaphore(0);
	PS[1] = new Semaphore(0);
}

// priority: 1 = high, 0 = low
void allocate(int priority) {
	M.wait();

	if (busy) {
		numWaiting[priority]++;
		M.notify();

		PS[priority].wait();

		// Re-enter critical section before mutating shared state.
		M.wait();
		numWaiting[priority]--;
	}

	busy = true;
	M.notify();
}

void deallocate() {
	M.wait();
	busy = false;

	// Wake higher priority waiters first, then lower priority.
	if (numWaiting[1] > 0) {
		PS[1].notify();
	} else if (numWaiting[0] > 0) {
		PS[0].notify();
	}

	M.notify();
}

import core.thread;
import std.stdio;
import core.atomic;

shared int inCritical = 0;
shared int violations = 0;

void worker(int prio, int id) {
    allocate(prio);
    if (atomicOp!"+="(inCritical, 1) > 1) {
        atomicOp!"+="(violations, 1);
    }
    Thread.sleep(dur!"msecs"(100));
    atomicOp!"-="(inCritical, 1);
    deallocate();
    writeln("done ", id, " prio=", prio);
}

unittest {
    auto t1 = new Thread({ worker(0, 1); });
    auto t2 = new Thread({ worker(1, 2); });
    auto t3 = new Thread({ worker(1, 3); });

    t1.start(); t2.start(); t3.start();
    t1.join(); t2.join(); t3.join();

    assert(violations == 0, "More than one thread entered critical section");
}