package box

import (
	"time"
)

type seqJob struct {
	jobs     []Job
	interval time.Duration
}

func (s seqJob) Run() {
	for _, j := range s.jobs {
		j.Run()
		time.Sleep(s.interval)
	}
}

func SeqJobs(interval time.Duration, jobs ...Job) Job {
	return seqJob{
		jobs:     jobs,
		interval: interval,
	}
}

type Loop struct {
	job      Job
	interval time.Duration
}

func NewLoop(interval time.Duration, j Job) Loop {
	return Loop{
		job:      j,
		interval: interval,
	}
}

func (l Loop) loop() {
	for {
		l.job.Run()
		time.Sleep(l.interval)
	}
}

func (l Loop) Run() {
	go l.loop()
}
