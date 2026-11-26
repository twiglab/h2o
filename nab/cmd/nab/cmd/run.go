package cmd

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	defer cliMgr.Close()

	interval := taskInterval()
	onofflog := onoffLog()

	mqttLog()

	p := pool()
	defer p.Release()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	agent := &nab.Agent{
		Global:      g,
		IDB:         db,
		CliMgr:      cliMgr,
		Logger:      logger,
		OnoffLogger: onofflog,
		Pool:        p,
	}

	mcli, act := mqttcli(
		nab.ClientID(g.Box, nab.NAB),
		nab.SubscriptTopic(g.Box),
		nab.AgentHandle(agent),
	)
	defer mcli.Disconnect(500)

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
	if !token.WaitTimeout(60 * time.Second) {
		slog.ErrorContext(ctx, "mqtt connect timeout")
		return errors.New("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		slog.ErrorContext(ctx, "mqtt connect error", slog.Any("error", err))
		return err
	}

	lps.Run()

	mux := chi.NewMux()
	mux.Mount("/debug", middleware.Profiler())
	mux.HandleFunc("/version", nab.VersionHandle)

	srv := httpserver(mux)

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			// 直接停止进程
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	log.Println("quit...")
	return nil
}
