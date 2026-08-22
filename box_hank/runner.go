package box

import "github.com/twiglab/h2o/box/internal/cron"

func NewCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

func Seq() error {
	return nil
}
