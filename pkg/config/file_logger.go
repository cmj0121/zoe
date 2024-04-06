package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cmj0121/zoe/pkg/service/types"
)

// Save the logger and the writer to the persistent file
type FileLogger struct {
	*os.File
}

func NewFileLogger(path string) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	logger := &FileLogger{
		File: file,
	}

	return logger, nil
}

func (f *FileLogger) Write(msg *types.Message) error {
	switch data, err := json.Marshal(msg); err {
	case nil:
		data := fmt.Sprintf("%s\n", data)
		_, err := f.WriteString(data)
		return err
	default:
		return err
	}
}

func (f *FileLogger) Close() error {
	return f.Close()
}
