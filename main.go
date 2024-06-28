package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker 			run image <cmd> <params>
// go run main.go	run 	  <cmd> <params>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()

	default:
		panic("bad command")
	}
}

func run(){
	fmt.Printf("Running %v\n", os.Args[2:])

	// /proc/self/exe will ensure it run itself again
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// creating some namespaces, using system process attributes
	// syscall.SysProcAttr is a struct that allows you to specify operating system-specific attributes for new processes. It is used in conjunction with the os/exec package to customize how processes are created and managed.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// cloning is what creates a new process, we can arbitrary command in.
		// we gonna start with flag for unix time sharing system.
		// All there in UTS namespace is hostname. 
		// But this is gonna let us have our own hostname inside the container, so can see its own and can't see what's happening on the host
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Run()
}

func child(){
	fmt.Printf("Running child %v\n", os.Args[2:])
	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// To run: sudo go run main.go run /bin/bash