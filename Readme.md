
Code is from this [Talk](https://www.youtube.com/watch?v=8fi7uSYlOdc&ab_channel=GOTOConferences) from GoTo 2018


## Namespaces:

- It is a concept used in computing to create isolated environments.

- Containers use namespaces to create isolated environments for running applications.

- Namespaces is where we limit what a process can see.

- Created with syscalls

- This is a big part of a container, it makes a container what it is, restricting the view the processes have of the host machine. 

Following are different types of namespaces:

    - Unix Timesharing System 
    - Process IDs 
    - Mounts 
    - Network 
    - User IDs 
    - InterProcess Comms


## CGroups

- What you can use 
- Filesystem interface 
    - Memory 
    - CPU 
    - I/O 
    - Process numbers



## Working of container:

1. Encapsulation into a container:

The entity that encapsulates these namespaces into what we call a "container" is not a single technological component, but rather the container runtime (like Docker). The runtime creates and manages these namespaces together, along with other isolation features like cgroups.

2. Functioning inside a container:

While each namespace isolates a particular aspect of the system, they work together within the container to provide a complete, isolated environment. Here's how:

- `UTS namespace`: Isolates hostname and domain name

- `PID namespace`: Provides an isolated process tree

- `Network namespace`: Isolates network interfaces, routing tables, etc.

- `Mount namespace`: Provides an isolated file system view

- `IPC namespace`: Isolates inter-process communication resources

- `User namespace`: Isolates user and group ID number spaces

These namespaces aren't completely separate; they interact with each other within the confines of the container. 

3. Container Creation Process:

    When a container is created:

    - The container runtime creates new instances of each required namespace

    - It then launches the container's init process (e.g., your application) within these namespaces

    - This init process becomes PID 1 in the container's `PID namespace`