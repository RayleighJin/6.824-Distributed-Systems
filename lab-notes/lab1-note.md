# Lab Process

![lab1-process](/Users/hao/Personal Projects/6.824/lab-notes/lab1-process.jpg)

# Program Design

## Object Design

### Master

* mapTask, a list of Task
* reduceTask, a list of Task
* intermidiateFile, a list of map(string -> )
* nReduce, int
* masterState
* complete, bool

### Task

* taskType, int, 0 is map while 1 is reduce
* taskID, corresponding to intermediate file
* filename
* taskState
* nReduce
* nFile
* time

## Function Design

### Master.go

* MakeMaster
* server, register a thread in rpc
* AssignTask, assign task to worker
* UpdateTask, update the task status in Master
* UpdateMaster, map -> reduce -> complete

### Worker.go

* Worker, infinite loop, keep sending request
* MapTask
* ReduceTask