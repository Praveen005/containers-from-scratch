
code is from this [Talk](https://www.youtube.com/watch?v=8fi7uSYlOdc&ab_channel=GOTOConferences) from GoTo 2018


## Namespaces:

- Namespaces is where we limit what a process can see.
- Created with syscalls
- This is a big part of a container, it makes a container what it is, restricting the view the processes have of the host machine. 
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