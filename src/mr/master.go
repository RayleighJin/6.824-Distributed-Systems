package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

// Master
type Master struct {
	// Your definitions here.
	Initialized   bool
	activeWorkers map[string]*Worker
	todoTaks      []Task
	mapWorkers    int
	reduceWorkers int
	mapTasks      int
	reduceTasks   int
	nRuduce       int
	mu            sync.RWMutex
}

// Task is exported
type Task struct {
	// Task is either a map task or reduce task
	ID    int
	state int
}

// Your code here -- RPC handlers for the worker to call.

// Init func
func (m *Master) Init(files []string, nReduce int) {
	m.activeWorkers = map[string]*Worker{}
	m.todoTaks = []Task{}
	m.nRuduce = nReduce // number of reduce tasks
	m.InitStatus()
	m.CreateLocalFiles(files, nReduce)
	m.Initialized = true
}

// InitStatus func
func (m *Master) InitStatus() {
	m.mu.Lock()
	m.mapWorkers = 0
	m.reduceWorkers = 0
	m.mapTasks = 0
	m.reduceTasks = 0
	m.mu.Unlock()
}

// CreateLocalFiles creates files
func (m *Master) CreateLocalFiles(files []string, nReduce int) {
	var tmpMaps []string
	for i := 0; i != nReduce; i++ {
		file := "mr-out-" + strconv.Itoa(i+1)
		tmpMap := "mr-tmp-map-" + strconv.Itoa(i+1)
		tmpMaps = append(tmpMaps, tmpMap)
		if f, e := os.Create(file); e != nil {
			log.Printf("[Master] creates [%s] file successfully", file)
		}
		if f, e := os.Create(tmpMap); e != nil {
			log.Printf("[Master] creates [%s] file successfully", tmpMap)
		}
	}
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
// //
// func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
// 	reply.Y = args.X + 1
// 	return nil
// }

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.

	m.server()
	return &m
}
