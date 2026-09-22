package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/twiglab/h2o/chrgg"
	"github.com/twiglab/h2o/chrgg/gql"
	"github.com/twiglab/h2o/chrgg/orm"
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

	sl := serverLog()

	mcli := mqttcli()
	entc := entcli()

	act := chrgg.NewMQTTAction(mcli)
	db := &orm.DBx{Cli: entc}

	svr := &chrgg.ChargeServer{
		Sender: act,
		DBx:    db,
		Logger: sl,
		WAL:    cwal(),
		MCli:   mcli,
	}

	if err := svr.Run(); err != nil {
		log.Fatal(err)
	}

	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Mount("/gql", gql.Handle(db))

	return http.ListenAndServe(webaddr(), mux)
}
