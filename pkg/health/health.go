package health

/*
Key Differences Between Livez and Readyz
Liveness (/livez):
    Determines if the application is running without deadlocks or critical failures
    Failure triggers container restart in Kubernetes
    Should only fail if the application needs to be restarted
Readiness (/readyz):
    Determines if the application is ready to serve requests
    Failure removes the pod from load balancer rotation
    Can fail temporarily during startup, configuration updates, or dependency issues

What I’ve described in this section is a polling approach, but I’d like to point out
that control loops that broadcast heartbeats are another way of implementing this pat-
tern. Using this technique, components are continually broadcasting their lifecycle
state with one control loop, and entities that have an interest in this app’s status will be
listening for these events and responding appropriately. I’ll remind you of the discus-
sion of chapter 4: what I’m describing here is an event-driven pattern. As the archi-
tect/developer, you must understand the architectural patterns of the software as a
whole and design/implement appropriately. The key with either approach is that con-
trol loops provide the redundancy that compensates for the uncertainty inherent in
distributed systems.

*/
