package nab

import (
	"cmp"
	"slices"
	"time"

	"github.com/twiglab/h2o/nab/orm/ent"
)

func SplitByCli(list []*ent.Dev) [][]*ent.Dev {
	slices.SortFunc(list, func(a, b *ent.Dev) int { return cmp.Compare(a.Cli, b.Cli) })
	returnData := make([][]*ent.Dev, 0)
	i := 0
	var j int
	for {
		if i >= len(list) {
			break
		}
		for j = i + 1; j < len(list) && list[i].Cli == list[j].Cli; j++ {
		}

		returnData = append(returnData, list[i:j])
		i = j
	}
	return returnData
}

type Loops struct {
	jobs []Job
}

func NewLoops() Loops {
	return Loops{
		jobs: make([]Job, 0),
	}
}

func (l *Loops) AddJob(j Job) *Loops {
	l.jobs = append(l.jobs, j)
	return l
}

func (l *Loops) AddToNewLoop(delay time.Duration, job Job) *Loops {
	return l.AddJob(NewLoop(delay, job))
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
