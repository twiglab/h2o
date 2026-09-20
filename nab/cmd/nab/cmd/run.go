package cmd

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"context"
	"slices"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/gql"
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

	idb := db()
	defer idb.Close()

	ctx := context.Background()

	g := global(idb)
	mcli := mqttcli()
	act := sender(mcli)

	lps := nab.NewLoops()

	devs, err := idb.AllDev(ctx)
	if err != nil {
		return err
	}

	for s := range slices.Chunk(devs, 10) {
		t := nab.CollectTask{
			Global: g,
			Data:   s,
			Sender: act,
			Logger: logger,
			Delay:  1 * time.Second,
		}
		lps.AddToNewLoop(1*time.Second, t)
	}

	log.Println("run after 5s")
	time.Sleep(5 * time.Second)

	lps.Run()

	hd := nab.HandleData{
		Global: g,
		Sender: act,
		Logger: logger,
	}

	t := mcli.Subscribe(nab.SubscriptTopic(g.Box), 0x01, nab.OnOffHandle(hd))
	t.Wait()

	if err := t.Error(); err != nil {
		return err
	}

	mux := chi.NewMux()
	mux.Mount("/gql", gql.Handle(hd, gql.WithPath("/gql")))

	return http.ListenAndServe(webaddr(), mux)
}
