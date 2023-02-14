package main

import (
	"github.com/sirupsen/logrus"
	"github.com/trusted-defi/trusted-engine/cmd/trustedengine/version"
	"github.com/trusted-defi/trusted-engine/node"
	"github.com/trusted-defi/trusted-engine/service"
	"github.com/trusted-defi/trusted-engine/tools"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

var log = logrus.WithField("prefix", "main")

func main() {
	app := cli.App{}
	app.Name = "trustedengine"
	app.Usage = "this is a txpool runing in enclave"
	app.Action = startNode
	app.Version = version.Version()
	app.Commands = []*cli.Command{}
	//app.Flags = appFlags

	app.Before = func(ctx *cli.Context) error {
		// init log
		formatter := new(tools.TextFormatter)
		formatter.TimestampFormat = "2006-01-02 15:04:05"
		formatter.FullTimestamp = true
		formatter.DisableColors = true
		logrus.SetFormatter(formatter)
		logrus.SetLevel(logrus.TraceLevel)

		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}

	defer func() {
		if x := recover(); x != nil {
			panic(x)
		}
	}()

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

func startNode(ctx *cli.Context) error {
	log.Info("start node")
	n := node.NewNode()
	service.StartTrustedService(n)
	return nil
}
