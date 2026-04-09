import core.sync.condition;
import core.sync.mutex;
import core.thread;
import core.time;
import std.algorithm;
import std.conv;
import std.exception;
import std.stdio;

// Shared monitor state
__gshared Mutex mtx;
__gshared Condition cv;
__gshared bool busy = false;

struct Waiter {
    int priority; // 1 = high, 0 = low
    long ticket;  // FIFO tie-break within same priority
}

__gshared Waiter[] waitQueue;
__gshared long nextTicket = 0;

shared static this() {
    mtx = new Mutex();
    cv = new Condition(cast(Mutex) mtx);
}

private bool isMyTurn(int priority, long ticket) {
    if (waitQueue.length == 0) {
        return false;
    }

    // Queue is kept sorted: higher priority first, then FIFO by ticket.
    return waitQueue[0].priority == priority && waitQueue[0].ticket == ticket;
}

void allocate(int priority) {
    enforce(priority == 0 || priority == 1, "priority must be 0 or 1");

    mtx.lock();
    scope(exit) mtx.unlock();

    const myTicket = nextTicket++;
    waitQueue ~= Waiter(priority, myTicket);

    // Priority queue + stable order among same-priority waiters.
    sort!((a, b) => (a.priority > b.priority) ||
        (a.priority == b.priority && a.ticket < b.ticket))(waitQueue);

    while (busy || !isMyTurn(priority, myTicket)) {
        cv.wait(); // Temporarily unlocks mtx and re-locks before returning.
    }

    // We are the selected waiter; remove ourselves from the queue now.
    waitQueue = waitQueue[1 .. $];
    busy = true;
}

void deallocate() {
    mtx.lock();
    scope(exit) mtx.unlock();

    busy = false;

    cv.notifyAll();
}

// ----------------------------
// Simple tests
// ----------------------------

__gshared int inCritical = 0;
__gshared int violations = 0;
__gshared int[] acquireOrder;
__gshared Mutex testLogMtx;

shared static this() {
    if (testLogMtx is null) {
        testLogMtx = new Mutex();
    }
}

void worker(int id, int priority, Duration holdFor) {
    allocate(priority);

    testLogMtx.lock();
    inCritical += 1;
    if (inCritical > 1) {
        violations += 1;
    }
    acquireOrder ~= id;
    testLogMtx.unlock();

    Thread.sleep(holdFor);

    testLogMtx.lock();
    inCritical -= 1;
    testLogMtx.unlock();

    deallocate();
}

unittest {
    // Reset test state.
    mtx.lock();
    busy = false;
    waitQueue.length = 0;
    nextTicket = 0;
    mtx.unlock();

    inCritical = 0;
    violations = 0;
    acquireOrder.length = 0;

    // Start low first so it holds the resource.
    auto low1 = new Thread({ worker(1, 0, dur!"msecs"(120)); });
    low1.start();
    Thread.sleep(dur!"msecs"(20));

    // These queue while low1 is in critical section.
    auto low2 = new Thread({ worker(2, 0, dur!"msecs"(40)); });
    auto high = new Thread({ worker(3, 1, dur!"msecs"(40)); });
    low2.start();
    Thread.sleep(dur!"msecs"(5));
    high.start();

    low1.join();
    low2.join();
    high.join();

    assert(violations == 0, "More than one thread entered critical section");

    // Expected acquisition order:
    // low1 first (it started first and grabbed resource),
    // then high should win over queued low2.
    assert(acquireOrder.length == 3);
    assert(acquireOrder[0] == 1);
    assert(acquireOrder[1] == 3, "High priority should acquire before low2");
}

