package main

import (
	"github.com/simabdi/vodka-authservice/config"
)

func main() {
	config.Initialize()
	config.Connection()
}
