package main

import (
	"github.com/trusted-defi/trusted-engine/cmd/trustedengine/version"
	"github.com/trusted-defi/trusted-engine/log"
	"github.com/trusted-defi/trusted-engine/node"
	"github.com/trusted-defi/trusted-engine/service"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

func main() {
	app := cli.App{}
	app.Name = "trustedengine"
	app.Usage = "this is a txpool runing in enclave"
	app.Action = startNode
	app.Version = version.Version()
	app.Commands = []*cli.Command{}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "generate",
			Value: false,
			Usage: "generate test secretdb",
		},
	}
	//app.Flags = appFlags

	app.Before = func(ctx *cli.Context) error {
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
	// init log
	log.InitLog()

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Info("start node")
	n := node.NewNode(ctx.Bool("generate"))
	service.StartTrustedService(n)
	return nil
}
