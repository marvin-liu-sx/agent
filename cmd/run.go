package cmd

import (
	"bee-agent/proxy"
	"bee-agent/utils"
	"context"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run a beekeeper agent",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "server",
			Usage:       "server url",
			DefaultText: "https://beekeeper.one",
			EnvVars:     []string{"BEE_AGENT_SERVER_URL"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "server authorization token",
			EnvVars: []string{"BEE_AGENT_SERVER_TOKEN"},
		},
		&cli.StringFlag{
			Name:        "api",
			Usage:       "local bee node api url",
			DefaultText: "http://localhost:1633",
			EnvVars:     []string{"BEE_AGENT_LOCAL_API"},
		},
		&cli.StringFlag{
			Name:        "debug-api",
			Usage:       "local bee node debug api url",
			DefaultText: "http://localhost:1635",
			EnvVars:     []string{"BEE_AGENT_LOCAL_DEBUG_API"},
		},
		&cli.Int64Flag{
			Name:        "port",
			Usage:       "listen port",
			DefaultText: "50505",
			EnvVars:     []string{"BEE_AGENT_LISTEN_PORT"},
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "debug",
		},
		&cli.BoolFlag{
			Name:  "disable-register",
			Usage: "disable register to server",
		},
	},
	Action: func(cctx *cli.Context) error {

		if cctx.Bool("debug") {
			_ = logging.SetLogLevel("*", "debug")
		} else {
			_ = logging.SetLogLevel("*", "info")
		}

		server := cctx.String("server")
		if server == "" {
			server = "https://beekeeper.one"
		}

		api := cctx.String("api")
		if api == "" {
			api = "http://localhost:1633"
		}
		debugAPI := cctx.String("debug-api")
		if debugAPI == "" {
			debugAPI = "http://localhost:1635"
		}

		port := cctx.Int64("port")
		if port == 0 {
			port = 50505
		}

		m := proxy.NewManager(api, debugAPI, cctx.String("server"), cctx.String("token"), port, cctx.Bool("disable-register"))

		ctx := utils.ReqContext(cctx)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			if err := m.Start(); err != nil {
				fmt.Printf("start err: %s", err)
			}
			cancel()
		}()

		<-ctx.Done()
		m.Stop()
		return nil
	},
}
