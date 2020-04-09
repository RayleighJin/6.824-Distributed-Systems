package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"time"
)

// KeyValue is the k/v pair struct
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// Worker is the called by mrworker.go
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the master.
	// CallExample()
	for {
		task := AssignTask()
		switch task.Flag {
		case 1:
			// all tasks are being taken care of, wait for another sec to send request
			time.Sleep(time.Second)
			continue
		case 2:
			// all work done
			return
		}
		switch task.Task.taskType {
		case mapTask:
			MapTask(mapf, task.Task)
			UpdateTask(task.Task.taskType, task.Task.taskID)
		case reduceTask:
			MapTask(mapf, task.Task)
			UpdateTask(task.Task.taskType, task.Task.taskID)
		default:
			panic("worker panic")
		}
		time.Sleep(time.Second)
	}

}

// MapTask executes the mapf
func MapTask(mapf func(string, string) []KeyValue, task Task) {

}

// ReduceTask executes the reducef
func ReduceTask(reducef func(string, string) []KeyValue, task Task) {

}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
// func CallExample() {

// 	// declare an argument structure.
// 	args := ExampleArgs{}

// 	// fill in the argument(s).
// 	args.X = 99

// 	// declare a reply structure.
// 	reply := ExampleReply{}

// 	// send the RPC request, wait for the reply.
// 	call("Master.Example", &args, &reply)

// 	// reply.Y should be 100.
// 	fmt.Printf("reply.Y %v\n", reply.Y)
// }

// AssignTask gets the task information from master
func AssignTask() AssignTaskReply {
	args := ExampleArgs{}
	reply := AssignTaskReply{}
	call("Master.AssignTask", &args, &reply)
	return reply
}

// UpdateTask updates the state of tasks in master
func UpdateTask(taskType int, taskID int) {
	// TODO
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
