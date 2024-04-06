package config

import (
	"fmt"

	"github.com/cmj0121/zoe/pkg/service/types"
)

// Save the logger and the writer to the console
type ConsoleLogger struct{}

func (c *ConsoleLogger) Write(msg *types.Message) error {
	fmt.Println(msg.String())
	return nil
}

func (c *ConsoleLogger) Close() error {
	return nil
}
