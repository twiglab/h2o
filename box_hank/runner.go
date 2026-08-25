package box

import (
	"github.com/twiglab/h2o/box/internal/cron"
)

type Job interface {
	Run()
}

type CronExec struct {
	cron      *cron.Cron
	isRunning bool
}

func NewCronExec() *CronExec {
	return &CronExec{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (c *CronExec) AddJob(spec string, j Job) *CronExec {
	if c.isRunning {
		panic("is running")
	}
	c.cron.AddJob(spec, j)
	return c
}

func (c *CronExec) Run() {
	c.cron.Start()
	c.isRunning = true
}
