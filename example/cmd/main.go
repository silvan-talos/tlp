package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/silvan-talos/tlp/example/http"
	"github.com/silvan-talos/tlp/example/mysql"
	"github.com/silvan-talos/tlp/example/user"
	"github.com/silvan-talos/tlp/log"
	"github.com/silvan-talos/tlp/logging"
)

var logLevel = &cli.StringFlag{
	Name:  "log-level",
	Usage: "Specifies the log level. Options: error | warn | info | debug",
	Value: "info",
}

func main() {
	app := newCLIApp()
	app.Action = startServer
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func newCLIApp() *cli.App {
	app := cli.NewApp()
	app.Name = "example"
	app.Usage = "testing logging library"
	app.Flags = []cli.Flag{
		logLevel,
	}
	return app
}

func startServer(cliCtx *cli.Context) error {
	logLevelFlag := cliCtx.String(logLevel.Name)
	logger, err := log.Default().WithAttrs(
		logging.NewAttr("server", "example"),
		logging.NewAttr("env", "dev"),
	).WithLevel(logLevelFlag)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}
	log.SetDefault(logger)

	exitChan := make(chan error, 1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
		exitChan <- fmt.Errorf("signal: %s", <-quit)
	}()

	users := mysql.NewUserRepository()
	userService := user.NewService(users)
	server := http.NewServer(userService)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("create server listener: %w", err)
	}
	go func() {
		err = server.Serve(lis)
		exitChan <- fmt.Errorf("server failed: %w", err)
	}()

	log.Info(cliCtx.Context, "server stopped", "reason", <-exitChan)
	return nil
}
