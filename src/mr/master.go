package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

// state of the master
const (
	newMaster = iota
	completeMap
	completeReduce
)

// task type
const (
	mapTask = iota
	reduceTask
)

// state of a task
const (
	initialized = iota
	inProgress
	complete
)

// Master is the master obj
type Master struct {
	// Your definitions here.
	mapTask          []Task
	reduceTask       []Task
	intermediateFile map[int][]string // file name -> taskid
	nReduce          int
	masterState      int
	complete         bool
}

// Task is exported
type Task struct {
	// Task is either a map task or reduce task
	taskType  int
	taskID    int
	filename  string
	taskState int
	nReduce   int
	nFile     int
	time      time.Time
}

// AssignTask assigns task to worker by master
// Your code here -- RPC handlers for the worker to call.
func (m *Master) AssignTask(_ *ExampleArgs, reply *AssignTaskReply) error {
	if m.masterState == newMaster {
		// master is just created, no work has been distributed
		// assign map task to workers
		for i, task := range m.mapTask {
			if task.taskState == initialized ||
				task.taskState == inProgress &&
					time.Now().Sub(m.mapTask[i].time) > time.Duration(10)*time.Second {
				// task is just initialized and not processed yet OR
				// task is processed but excceeds the time limit
				// assign the map task
				reply.Task.taskType = task.taskType
				reply.Task.taskID = task.taskID
				reply.Task.filename = task.filename
				reply.Task.taskState = task.taskState
				reply.Task.nReduce = task.nReduce
				reply.Flag = 0

				m.mapTask[i].taskState = inProgress
				m.mapTask[i].time = time.Now()
				return nil
			}
		}
		reply.Flag = 1
	} else if m.masterState == completeMap {
		// map tasks all done, now attribute the reduce task
		for i, task := range m.reduceTask {
			if task.taskState == initialized ||
				task.taskState == inProgress &&
					time.Now().Sub(m.mapTask[i].time) > time.Duration(10)*time.Second {
				// task is initialized and not processed yet OR
				// task is processed but excceeds the time limit
				// assign the reduce task
				reply.Task.taskType = task.taskType
				reply.Task.taskID = task.taskID
				reply.Task.filename = task.filename
				reply.Task.taskState = task.taskState
				reply.Task.nReduce = task.nReduce
				reply.Flag = 0

				m.reduceTask[i].taskState = inProgress
				m.reduceTask[i].time = time.Now()
				return nil
			}
		}
		reply.Flag = 1
	} else if m.masterState == completeReduce {
		// all work done
		reply.Flag = 2
	} else {
		log.Println("Something went wrong with masterState.")
	}
	return nil
}

// UpdateTask updates the task state in Master: initialized -> inProgress -> complete
func (m *Master) UpdateTask(args *UpdateTaskArgs, _ *ExampleReply) error {
	if args.taskType == mapTask {
		for i, task := range m.mapTask {
			if task.taskID == args.taskID {
				m.mapTask[i].taskState = complete
			}
		}
	} else if args.taskType == reduceTask {
		for i, task := range m.reduceTask {
			if task.taskID == args.taskID {
				m.reduceTask[i].taskState = complete
			}
		}
	}
	m.UpdateMaster()
	return nil
}

// UpdateMaster checks if jobs are done
// if all map jobs done, newMaster -> completeMap, and create reduce tasks
// if all reduce jobs done, completeMap -> completeReduce
func (m *Master) UpdateMaster() {
	if m.masterState == newMaster {
		// check if all map jobs done
		for _, currtask := range m.mapTask {
			if currtask.taskState != completeMap {
				return
			}
		}
		m.masterState = completeMap
		// create reduce tasks
		for i := 0; i < m.nReduce; i++ {
			m.reduceTask = append(m.reduceTask, Task{
				taskType:  reduceTask,
				taskID:    i,
				filename:  "",
				taskState: initialized,
				nReduce:   m.nReduce,
				nFile:     len(m.mapTask),
			})
		}
	} else if m.masterState == completeMap {
		for _, currtask := range m.reduceTask {
			if currtask.taskState != completeReduce {
				return
			}
		}
		m.masterState = completeReduce
		m.complete = true
	} else if m.masterState == completeReduce {
		m.complete = true
	} else {
		log.Println("Something went wrong with the masterState, please check")
	}
	return
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

// Done checks if the master has finished its job
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := m.complete
	return ret
}

// MakeMaster creates a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{
		nReduce:     nReduce,
		mapTask:     []Task{},
		masterState: newMaster,
		complete:    false,
	}
	for i, file := range files {
		m.mapTask = append(m.mapTask, Task{
			taskType:  mapTask,
			taskID:    i,
			filename:  file,
			taskState: initialized,
			nReduce:   m.nReduce,
		})
	}

	go m.server()
	return &m
}
