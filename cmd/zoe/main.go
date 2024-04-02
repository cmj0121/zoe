package main

import (
	"os"

	"github.com/cmj0121/zoe"
)

func main() {
	agent := zoe.New()
	os.Exit(agent.ParseAndRun())
}
