package handler

import (
	"MultiJobTimeTicker/config"
	"fmt"
	"runtime/debug"
	"sync"
)

type JobManager struct {
	Job            chan *config.JobContent
	wg             sync.WaitGroup
	closeSignal    chan interface{}
	completeSignal chan interface{}
}

var JOBManager *JobManager

func NewJobManager() *JobManager {
	var jobManger JobManager
	jobManger.Job = make(chan *config.JobContent)
	jobManger.closeSignal = make(chan interface{})
	jobManger.completeSignal = make(chan interface{})
	return &jobManger
}

func (jobManager *JobManager) RecvJob(job *config.JobContent) {
	jobManager.Job <- job
}

func (jobManager *JobManager) Start() {
	jobManager.wg.Add(1)
	go jobManager.dealJob()
}

func (jobManger *JobManager) dealJob() {
	defer  jobManger.wg.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(debug.Stack())
		}
	}()
	for {
		select {
		case job := <- jobManger.Job:
			go DoJob()
		case <- jobManger.closeSignal:
			break
		}
	}
}

func (jobManger *JobManager) GetCompleteSignal() <-chan interface{} {
	return jobManger.completeSignal
}

func (jobManager *JobManager) Close() {
	close(jobManager.closeSignal)
	go func () {
		jobManager.wg.Wait()
		close(jobManager.completeSignal)
	}()
}