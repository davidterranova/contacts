package main

import (
	"github.com/rs/zerolog/log"

	"github.com/davidterranova/contacts/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("failed to start contacts")
	}
}
