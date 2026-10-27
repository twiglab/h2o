package wecom

import (
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twiglab/h2o/dbcli/ent"
)

type reportCmd struct {
	tmpl *template.Template
	cli  *ent.Client

	Items []reportItem
	Now   time.Time
}

func (c *reportCmd) Do(ctx context.Context) error {
	c.Now = time.Now()
	rows, err := c.cli.QueryContext(ctx, report_sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	return sqlx.StructScan(rows, &c.Items)
}

func (c reportCmd) ToString() string {
	var sb strings.Builder
	if err := c.tmpl.ExecuteTemplate(&sb, template_report, c); err != nil {
		return err.Error()
	}
	return sb.String()
}

func reportCmdFn(cfg Global, args ...string) (Commander, error) {
	u := &reportCmd{
		tmpl: cfg.Template,
		cli:  cfg.Client,
	}
	return u, nil
}

type reportItem struct {
	Code string `db:"code"`

	DeviceCode string    `db:"device_code"`
	DeviceType string    `db:"device_type"`
	PosCode    string    `db:"pos_code"`
	Top        int64     `db:"top"`
	ChangeTime time.Time `db:"charge_time"`

	DataValue int64     `db:"data_value"`
	DataTime  time.Time `db:"data_time"`

	Diff int64 `db:"diff"`

	Rate    int64   `db:"rate"`
	PosName *string `db:"pos_name"`
	Alarm   int64   `db:"alarm"`

	DiffEng int64 `db:"diff_eng"`
}

const report_sql = `
select
	a.code,
	a.device_code,
	a.device_type,
	a.pos_code,
	a.top,
	a.charge_time
,
	b.dv as data_value,
	b.dt as data_time
,
	a.top - b.dv as diff
,
c.rate ,
c.pos_name,
c.alarm
,
(a.top - b.dv) as diff_eng

from
	(
	select
		code,
		device_code,
		device_type,
		pos_code,
		top,
		charge_time
	from
		(
		select
			*,
			row_number() over(partition by device_code order by charge_time desc) as num
		from
			device_top
) t
	where
		t.num = 1
) a
join
(
	select
		device_code,
		max(data_value) as dv ,
		max(data_time) as dt
	from
		nh_record
		where data_time >= CURRENT_DATE - INTERVAL '7 day'
	group by
		device_code
) b on
	(a.device_code = b.device_code)
join device c on
	(a.device_code = c.device_code)
order by c.pos_code
`
