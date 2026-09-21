package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"time"

	"github.com/spf13/cobra"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/gql"
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

	db := db()
	defer db.Close()

	g := global()
	mcli := mqttcli()
	act := sender(mcli)
	cliMgr := clientMgr(db)

	lps := loops(db, func(data []*ent.Dev) nab.Job {
		return nab.CollectTask{
			Global:    g,
			Data:      data,
			Sender:    act,
			IDB:       db,
			ClientMgr: cliMgr,
			Logger:    logger,
			Delay:     1 * time.Second,
		}
	})

	agent := &nab.Agent{
		Global: g,
		Sender: act,
		IDB:    db,
		CliMgr: cliMgr,
		MCli:   mcli,
		Logger: logger,
	}

	log.Println("run after 5s")
	time.Sleep(5 * time.Second)

	lps.Run()
	if err := agent.Start(); err != nil {
		log.Fatal(err)
	}

	root := chi.NewMux()
	root.Use(middleware.Recoverer, middleware.RequestID)
	root.Mount("/gql", gql.Handle(agent, gql.WithPath("/gql")))

	return http.ListenAndServe(webaddr(), root)
}
