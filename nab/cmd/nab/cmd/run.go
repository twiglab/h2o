package cmd

import (
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/twiglab/h2o/nab/box/driverbox"
	"github.com/twiglab/h2o/nab/box/exports"
	"github.com/twiglab/h2o/nab/box/plugins"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run() error {
	// 设置日志级别
	_ = os.Setenv("LOG_LEVEL", "debug")
	// 第一步: 启用内置Plugin
	// 可选: plugins.EnableAll() 启用所有内置Plugin
	// 或单独启用: 导入对应Plugin包,如 modbus.EnablePlugin()

	plugins.EnableAll()

	//modbus.EnablePlugin()
	//httpserver.EnablePlugin()

	// 第二步: 启用Export模块
	// 可选: exports.EnableAll() 启用所有内置Export模块
	// 或单独启用: 导入对应Export包,如 gateway.EnableExport()

	exports.EnableAll()
	// mirror.EnableExport()
	//mqtt.EnablePlugin()

	// 第三步: 启动 driver-box 服务
	// 1. 初始化环境配置
	// 2. 初始化日志记录器
	// 3. 启动所有 Export
	// 4. 启动所有 Plugin
	// 5. 触发服务状态事件
	err := driverbox.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer driverbox.Stop()
	// 第四步: 阻塞主线程
	// driver-box 服务会在后台运行
	return http.ListenAndServe(":10001", nil)
}
