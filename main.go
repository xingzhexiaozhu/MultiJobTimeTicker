package main

import (
	"MultiJobTimeTicker/config"
	"MultiJobTimeTicker/handler"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	conf := config.GetConfig()

	err := handler.Init(conf)
	if err != nil {
		log.Printf("Error %v\n", err)
		return
	}

	// start JobManager
	handler.JOBManager = handler.NewJobManager()
	handler.JOBManager.Start()

	// 启动定时任务
	crontabJob := cron.New()
	for jobID := 1; jobID < conf.GlobalConf.JobNum; jobID++ {
		jobIDStr := strconv.FormatInt(int64(jobID), 10)
		jobContent := conf.JobContent[jobIDStr]
		job := &config.Job{
			Type:      jobContent.Type,
			Title:     jobContent.Title,
			Content:   jobContent.Content,
			Url:       jobContent.Url,
			Freq:      jobContent.Freq,
			Condition: jobContent.Condition,
			Success:   0,
		}
		// 添加任务
		crontabJob.AddFunc(job.Freq, func() {
			handler.JOBManager.GenJob(job)
		})
	}
	crontabJob.Start()
	go signalProc()

	select {}
}

func signalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGUSR1)

	handler.JOBManager.Close()
	<- handler.JOBManager.GetCompleteSignal()
}