package cmd

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm/ent"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行采集服务",
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

	_ = rootLog()

	logger := serverLog()

	g := global()

	db := db(g)
	defer func() { _ = db.Close() }()

	cliMgr := clientMgr(db)
	interval := taskInterval()
	wl := cwal()

	agent := &nab.Agent{
		Global: g,
		IDB:    db,
		CliMgr: cliMgr,
		Logger: logger,
		WAL:    wl,
	}

	mcli, act := mqttcli2(nab.SubscriptTopic(g.Box), nab.AgentHandle(agent))

	lps := loops(db, func(data []*ent.Dev) nab.Job {
		return nab.CollectTask{
			Global:    g,
			Data:      data,
			Sender:    act,
			IDB:       db,
			ClientMgr: cliMgr,
			Logger:    logger,
			Delay:     interval,
		}
	})

	token := mcli.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatal("mqtt - connect error")
	}

	lps.Run()

	go http.ListenAndServe(webaddr(), nil)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("quit...")
	return nil
}
