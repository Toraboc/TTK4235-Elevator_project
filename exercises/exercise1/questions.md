Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is a about handling multiple tasks or changes happening at the same time. Paralellism is about preforming multiple tasks at the exact same time. feks. graphics computing across many GPU cores.

What is the difference between a *race condition* and a *data race*? 
>  A data race happens when two or more processes write to the memory at the exact same time, while a race condition happens more generally when something goes wrong because two or more tasks are preformed at the same time or in the wrong order.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> A scheduler keeps track of different tasks, and distributes the tasks among CPUs. It also keeps track of task runtime and suspends and continues tasks as they are blocked or priority.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> Multiple threads allow for multiple tasks to be preformed in paralell without having to wait for eachother to finish. Threads make sure the priority processes can be handeled without being sat on pause by slow processes that may take a while to finish.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> They can solve the same problems as normal threads but with less overhead. They are managed by language runtime and use less OS threads.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Both. Helps in the sense that it makes it easier to handle hard time critical tasks, but makes it harder in the sense of introducing race conditions and other hard to deal with things.

What do you think is best - *shared variables* or *message passing*?
> Shared variables are faster but more dangerous, when preformance is critical. Message passing is safer but slower.

