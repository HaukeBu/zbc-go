# TODO 


## feat/sub-callback

[x] Seperate managers into their own modules/services

[x] Refactor common stuff to common directory

[x] Refactor zbdispatch stuff to directory

[x] Refactor zbtopology

[x] Refactor zbtransport

[x] Refactor zbsubscribe

[x] Refactor zbsocket

[x] Improve error handling (example: when broker is not there)

[x] If you cannot connect to bootstrapAddr, return error

[x] Cleanup panics

[x] Use Partition information

[x] Write FailTask request

[x] Pass client to the task callback

[x] Pass client to the topic callback

[x] Improve RoundRobin controller

[x] Improve handling of credits on TaskSubscription

[x] Fix bug in Task/Topic subscription after refactoring 

[x] Improve RoundRobin test

[x] Improve subscriptions tests to create 1000+ workflow instances

[x] Test scoped/unscoped handlers on subscriptions

[x] Add timers for workflow instances

[x] Client should timeout if I kill broker (DATA RACE: problem)

[x] Replace dependency management (godep with dep)

[x] Fix TravisCI

[x] FIX: if the error happend in the middle of opening we should close it all

[x] Reimplement topic subscriptions to support easier testing

[x] Fix topologyTicker problem

[x] Prevent data race when broker crashes

...

[] Extension of credits should be easier to manipulate

[] Move closing of subscription to TaskSubscription object

[] Move closing of subscription to TopicSubscription object



[] Add benchmark for workflow instances

[] Test for error on TopicSubscription handler


