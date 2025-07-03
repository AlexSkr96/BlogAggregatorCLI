package main

import (
	"fmt"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
)

func main() {
	var state, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("error while getting config from file: %v\n", err)
	}
}
