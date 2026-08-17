package pick

import "github.com/twiglab/h2o/hank/pick/internal/cron"

func NewCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

func Seq() error {
	return nil
}
