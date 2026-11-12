package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"time"

	"github.com/spf13/cobra"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/orm/ent"
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

	_ = rootLog()

	logger := serverLog()

	g := global()

	db := db(g)
	defer db.Close()

	mcli := mqttcli()
	act := sender(mcli)
	cliMgr := clientMgr(db)

	interval := taskInterval()

	wl := cwal()

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

	agent := &nab.Agent{
		Global: g,
		Sender: act,
		IDB:    db,
		CliMgr: cliMgr,
		MCli:   mcli,
		Logger: logger,
		WAL:    wl,
	}

	log.Println("run after 5s")
	time.Sleep(5 * time.Second)

	lps.Run()
	if err := agent.Start(); err != nil {
		log.Fatal(err)
	}

	go http.ListenAndServe(webaddr(), nil)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	fmt.Println("quit...")
	return nil
}
