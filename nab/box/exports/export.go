package exports

import (
	"github.com/twiglab/h2o/nab/box/exports/discover"
	"github.com/twiglab/h2o/nab/box/exports/gateway"
	"github.com/twiglab/h2o/nab/box/exports/linkedge"
	"github.com/twiglab/h2o/nab/box/exports/mirror"
)

// EnableAll 加载driver-box框架内置的所有Export插件
// 功能:
//
//	依次调用各个内置Export的加载方法，包括基础Export、场景联动Export等
func EnableAll() {
	linkedge.EnableExport()
	mirror.EnableExport()
	discover.EnableExport()
	gateway.EnableExport()
}
