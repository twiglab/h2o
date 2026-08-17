package pick

import "github.com/twiglab/h2o/box/ocg/pick/internal/cron"

func NewCron() *cron.Cron {
	return cron.New(cron.WithSeconds())
}

func Seq() error {
	return nil
}
