package cmd

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/eyes"
	"github.com/twiglab/h2o/vigil/gql"

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
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run() {

	rootLog()

	cli := entcli()

	egg := eyes.NewElectricityEgg()
	hub := &vigil.Hub{
		ElectyMeterView: egg,
		Recorder:        vigil.WithRecorder(dbx(cli)),
		Logger:          serverLog(),
	}
	mcli := mqttcli()
	token := mcli.SubscribeMultiple(topics(), vigil.Handle(hub))
	token.Wait()

	gqlc := gql.NewConf(cli)

	mux := chi.NewMux()
	mux.Mount("/eyes", eyes.EyesMux(egg))
	mux.Mount("/gql", gql.Handle(gqlc))
	http.ListenAndServe(webaddr(), mux)
}
