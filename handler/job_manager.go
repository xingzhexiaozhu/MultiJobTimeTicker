package handler

import (
	"MultiJobTimeTicker/config"
	"log"
	"runtime/debug"
	"sync"
)

type JobManager struct {
	Job            chan *config.Job
	wg             sync.WaitGroup
	closeSignal    chan interface{}
	completeSignal chan interface{}
}

var JOBManager *JobManager

func NewJobManager() *JobManager {
	var jobManger JobManager
	jobManger.Job = make(chan *config.Job)
	jobManger.closeSignal = make(chan interface{})
	jobManger.completeSignal = make(chan interface{})
	return &jobManger
}

func (jobManager *JobManager) GenJob(job *config.Job) {
	jobManager.Job <- job
}

func (jobManager *JobManager) Start() {
	jobManager.wg.Add(1)
	go jobManager.dealJob()
}

func (jobManager *JobManager) dealJob() {
	defer  jobManager.wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error %v\n", debug.Stack())
		}
	}()
	for {
		select {
		case job := <- jobManager.Job:
			go MapUserData[job.Type].SelectUserAndDoJob(job)
		case <- jobManager.closeSignal:
			break
		}
	}
}

func (jobManager *JobManager) GetCompleteSignal() <-chan interface{} {
	return jobManager.completeSignal
}

func (jobManager *JobManager) Close() {
	close(jobManager.closeSignal)
	go func () {
		jobManager.wg.Wait()
		close(jobManager.completeSignal)
	}()
}