package main

import (
	"fmt"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/config"
)

func main() {
	myConfig, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error while getting config from file: %v\n", err)
	}
	err = myConfig.SetUser("Alex")
	if err != nil {
		fmt.Printf("error while setting user: %v\n", err)
	}
	myConfig, err = config.ReadConfig()
	if err != nil {
		fmt.Printf("error while getting config from file again: %v\n", err)
	}
	fmt.Printf("%v\n", myConfig)
}
