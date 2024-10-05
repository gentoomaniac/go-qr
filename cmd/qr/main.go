package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	gocli "github.com/gentoomaniac/go-qr/pkg/cli"
	"github.com/gentoomaniac/go-qr/pkg/logging"
	"github.com/gentoomaniac/go-qr/pkg/qr"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli struct {
	logging.LoggingConfig

	Data        string `help:"Data to be encoded" required:""`
	CodeVersion int    `help:"QR code version" default:"1"`

	Version gocli.VersionFlag `short:"V" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
		"commit":  commit,
		"binName": binName,
		"builtBy": builtBy,
		"date":    date,
	})
	logging.Setup(&cli.LoggingConfig)

	err, qrCode := qr.New(cli.CodeVersion, cli.Data)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	for _, line := range qrCode.ToString() {
		fmt.Println(line)
	}

	ctx.Exit(0)
}
