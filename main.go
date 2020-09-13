package main

import (
	"fmt"

	"github.com/gosidekick/gasgo/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("metadata: ", cfg.MetadataFile)
}
