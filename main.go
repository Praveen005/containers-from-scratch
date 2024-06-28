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
/*
PID TTY          TIME CMD
88447 pts/14   00:00:00 sudo
88448 pts/14   00:00:00 go
88567 pts/14   00:00:00 main
88573 pts/14   00:00:00 exe
88579 pts/14   00:00:00 bash
88635 pts/14   00:00:00 ps
*/
// we want the PIDs for the process inside the container to start from 1
// we can do it by using syscall.CLONE_NEWPID while creating the namespace
// Try running the command above again, and see
/*
  PID TTY          TIME CMD
98911 pts/14   00:00:00 sudo
98912 pts/14   00:00:00 go
99017 pts/14   00:00:00 main
99022 pts/14   00:00:00 exe
99027 pts/14   00:00:00 bash
99079 pts/14   00:00:00 ps
*/

// It still doesn't work, beacuse `ps` don't get the PIDs magically, but gets from the /proc directory, check by ls /proc
// so, inside my container it needs its own /proc directory, currently it is seeing the same /proc directory
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
	// A container should have its own root file system.
	// This is a fundamental principle of containerization, which ensures that containers are isolated from each other and the host system.
	// container uses as its root directory. It contains all the files and directories necessary for the container's processes to run.
	// Lizz used Vagarant-fs, we will use ubuntu file system.
	// way to do:
	// start the ubuntu container, extract its file system in a tar file, store it, make a root directory for your container, extract the fs here, and use it as the container's file system
	/*

	docker run -d --rm --name ubuntufs ubuntu:20.04 sleep 1000
	docker export ubuntufs -o ubuntufs.tar
	docker stop ubuntufs
	sudo mkdir -p /container-root
	sudo tar xf ubuntufs.tar -C /container-root/


	*/
	syscall.Chroot("/container-root")
	syscall.Chdir("/")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Start the container:  sudo go run main.go run /bin/bash
// In the container, start a process:  sleep 100   
// Go to another terminal outside the container i.e. from host: ps -C sleep
// get the PID(put in place of 55327), and run `sudo ls -l /proc/55327/root`
// you will get something like: lrwxrwxrwx 1 root root 0 Jun 28 18:48 /proc/55327/root -> /container-root
// you can see /container-root is the root directory of the running container