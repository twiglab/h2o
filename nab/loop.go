package nab

import (
	"time"
)

type Loops struct {
	jobs []Job
}

func (l Loops) Run() {
	for _, j := range l.jobs {
		j.Run()
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
