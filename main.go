package main

import (
	"context"

	"github.com/dvcrn/matrix-bridgekit/bridgekit"
	"github.com/dvcrn/matrix-bumble/bumble"
	"github.com/dvcrn/matrix-bumble/config"
)

func main() {
	env, err := bumble.Process()
	if err != nil {
		panic(err)
	}

	br := bridgekit.NewBridgeKit("Bumble", env.Localpart, "", "Bumble integration for Bumble", "1.0", &config.Config{}, config.ExampleConfig)
	connector := bumble.NewBumbleConnector(br)
	br.StartBridgeConnector(context.Background(), connector)
}
