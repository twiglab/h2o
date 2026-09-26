package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twiglab/h2o/dbcli"
	"github.com/twiglab/h2o/dbcli/ent"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showSql()
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)
}

func showSql() error {
	cli := entcli()
	defer cli.Close()

	return cli.Schema.WriteTo(context.Background(), os.Stdout)
}

func entcli() *ent.Client {
	name := viper.GetString("db.name")
	dsn := viper.GetString("db.dsn")

	cli, err := dbcli.OpenEntClient(name, dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("ent err: %w", err))
	}
	return cli
}
