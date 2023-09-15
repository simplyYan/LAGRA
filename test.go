package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	lagra "github.com/NervousGroove/LAGRA"
)

func BenchmarkLAGRA(b *testing.B) {
	b.SetParallelism(3)
	b.N = 14

	// Create a new LAGRA logger
	tomlConfig := `
        log_file = "app.log"
        log_level = "INFO"
    `
	logger, err := lagra.New(tomlConfig)
	if err != nil {
		fmt.Println("Error creating LAGRA logger:", err)
		return
	}
	for i := 0; i < b.N; i++ {
		logger.Info(nil, "This is a log message")
	}
}

func BenchmarkZerolog(b *testing.B) {
	// Create a new Zerolog logger
	zerolog.TimeFieldFormat = ""
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	for i := 0; i < b.N; i++ {
		log.Info().Msg("This is a log message")
	}
}
