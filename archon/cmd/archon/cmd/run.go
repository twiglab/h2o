package cmd

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/twiglab/h2o/archon/gql"

	"github.com/spf13/cobra"

	_ "net/http/pprof"
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

	http.Handle("/gql", playground.ApolloSandboxHandler("gql", "/gql/query"))
	http.Handle("/gql/query", gql.Handle(cli))

	if err := http.ListenAndServe(webaddr(), nil); err != nil {
		log.Fatal(err)
	}
}
