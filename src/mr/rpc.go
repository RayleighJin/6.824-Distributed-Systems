package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

// ExampleArgs is used as default
type ExampleArgs struct {
	X int
}

// ExampleReply is used as default
type ExampleReply struct {
	Y int
}

// AssignTaskReply is the reply by the method AssignTask
type AssignTaskReply struct {
	Flag int
	Task Task
}

// UpdateTaskArgs is used in UpdateTaskArgs
type UpdateTaskArgs struct {
	taskID   int
	taskType int
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
