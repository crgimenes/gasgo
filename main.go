package main

import (
	"fmt"

	"github.com/crgimenes/gasgo/config"
)

var (
	GitTag string = "dev"
)

func main() {
	config.GitTag = GitTag
	err := config.Load()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Version:", config.CFG.Version)
}
