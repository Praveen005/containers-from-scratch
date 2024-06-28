package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker 			run image <cmd> <params>
// go run main.go	run 	  <cmd> <params>

// we still don't have namespace process id.
// if you run: sudo go run main.go run /bin/bash and then run `ps`, it will show you large numbered PIDs
// we want the PIDs for the process inside the container to start from 1
// we can do it by using syscall.CLONE_NEWPID while creating the namespace
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
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	cmd.Run()
}

func child(){
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}