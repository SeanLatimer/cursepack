package main

import (
	"errors"
	"os"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/urfave/cli/v2"
)

var version = "0.1.0"
var commit = ""

const TEMPDIR_NAME = "cursepack"

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Value:   false,
				Usage:   "Enables verbose logging",
				EnvVars: []string{"CP_VERBOSE"},
			},
			&cli.BoolFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "Enables installing for a server, will download and install Forge automatically",
				EnvVars: []string{"CP_SERVER"},
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   "",
				Usage:   "Directory to install the pack in",
				EnvVars: []string{"CP_DIR"},
			},
		},
		Name:    "cursepack",
		Action:  run,
		Before:  before,
		Version: version + " " + commit,
	}

	err := app.Run(os.Args)
	if err != nil {
		jww.FATAL.Fatal(err)
	}
}

func before(ctx *cli.Context) error {
	var err error
	jww.SetStdoutThreshold(jww.LevelInfo)
	if ctx.Bool("verbose") {
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
	if ctx.String("dir") == "" {
		wd := ""
		wd, err = os.Getwd()
		ctx.Set("dir", wd)
	}
	return err
}

func run(ctx *cli.Context) error {
	pack := ctx.Args().Get(0)
	if os.Getenv("CP_PACK") != "" {
		pack = os.Getenv("CP_PACK")
	}
	if pack == "" {
		return errors.New("Must provide a pack")
	}
	if strings.HasSuffix(pack, ".zip") {
		return handleZipPack(PackInstallOptions{
			Pack:   pack,
			Path:   ctx.String("dir"),
			Server: ctx.Bool("server"),
		})
	}
	// IMPLEMENT
	return errors.New("Pack IDs are not currently supported")
}

type PackInstallOptions struct {
	Pack   string
	Server bool
	Path   string
}
