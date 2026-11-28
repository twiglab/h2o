package nab

import (
	"context"
)

type JobV2 interface {
	Run(ctx context.Context)
}

type Loops2 struct {
	loops []JobV2
}

func NewLoopsV2() *Loops2 {
	return &Loops2{loops: make([]JobV2, 0)}
}

func (l *Loops2) AddToNewLoop(job JobV2) *Loops2 {
	l.loops = append(l.loops, job)
	return l
}

// Run 启动全部循环。
func (l *Loops2) Run(ctx context.Context) {
	for _, lp := range l.loops {
		go lp.Run(ctx)
	}
}
