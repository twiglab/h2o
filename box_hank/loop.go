package box

import (
	"time"
)

type seqJob struct {
	jobs []Job
}

func (s seqJob) Run() {
	for _, j := range s.jobs {
		j.Run()
	}
}

func SeqJobs(jobs ...Job) Job {
	return seqJob{
		jobs: jobs,
	}
}

type Loop struct {
	job       Job
	isRunning bool
}

func NewLoop(j Job) Loop {
	return Loop{job: j}
}

func (l Loop) loop() {
	for {
		l.Run()
		time.Sleep(1 * time.Second)
	}
}

func (l Loop) Run() {
	l.isRunning = true
	go l.loop()
}
