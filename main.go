package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Run()
}

func child(){
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("container"))
	
	syscall.Chroot("/container-root")
	syscall.Chdir("/")

	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	syscall.Unmount("/proc", 0)
}

func cg() {
	cgroups := "/sys/fs/cgroup/"  // control group is going to this dir
	pids := filepath.Join(cgroups, "pids") // inside that to pids
	os.Mkdir(filepath.Join(pids, "praveen"), 0755) // inside that a cgroup I am creating called "praveen"
	must(ioutil.WriteFile(filepath.Join(pids, "praveen/pids.max"), []byte("20"), 0700)) // and we will write a limit of 20 on the number of processes, meaning there can be only 20 processes in my control group
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "praveen/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "praveen/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)) // basically getting current PIDs by os.Getpid() and writing inside my cgroup called `cgroup.procs`. Basically saying, this process is now the ember of this control group
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// If you go to `cd /sys/fs/cgroup` and `ls`, you will see directory of each of the type of control group that you can setup
// lets see memory: `cd memory` then `ls`: you will large number of parameters that you can set for memory
// In here you can also see docker folder, let's see inside: `cd docker` then `ls`, you can see lots of parameters.
// let's run a container: docker run --rm -it ubuntu /bin/bash
// Got a container id starting with c33718d4e517
// let's head back to another terminal and again do `ls` again inside docker directory, you will see a directory inside the control group structure starting with c33718d4e517, means docker has created a control group for this container.
// if you look for something inside that directory, by `cat c33718d4e517fb3bec3a00692e7faa92fd725fc6d85f124d70d0f61660a0c68e/memory.limit_in_bytes`
// got: 9223372036854771712  // probably using the max available. let's constraint it and see
// docker run --rm -it --memory=10M ubuntu /bin/bash
// `cat af54156ec22117c0ca6a699f1c5a579ae821b5e6961e5d26cf2d5aeecca51297/memory.limit_in_bytes` gave: 10485760
// what happened is docker wrote that number in that file and that how it tel the kernel to limit that particular container to 10MB memory

// Now comeback to /sys/fs/cgroup by `cd ../..`
// cd pids then cat docker/pids.max : it gives `max` means it will allow max number of processes to run inside the container
// Let's create a control group that limits the number of process that can run: see cg()

// Fault Bomb: :() { : | : &}; :, continously creates processes, don't run o host machine