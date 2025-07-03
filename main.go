package main

import (
	"fmt"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
)

func main() {
	var state config.State
	var err error
	state.Config, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("ERROR while getting config from file: %v\n", err)
	}
}
