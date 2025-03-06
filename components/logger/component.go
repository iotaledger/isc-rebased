package logger

import (
	"github.com/iotaledger/hive.go/app"
	"github.com/iotaledger/wasp/packages/evm/evmlogger"
)

func init() {
	Component = &app.Component{
		Name:      "Logger",
		Configure: configure,
	}
}

var Component *app.Component

func configure() error {
	log := Component.App().Logger.NewChildLogger("Ethereum")
	evmlogger.Init(log)
	return nil
}
