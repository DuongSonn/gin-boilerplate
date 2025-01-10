package logger

import (
	"log/slog"
	"os"
)

var (
	logger     *slog.Logger
	maskedKeys = []string{"password"}
)

func Init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			for _, key := range maskedKeys {
				if a.Key == key {
					a.Value = slog.StringValue("***")
				}
			}

			return a
		},
	}))
}

func GetLogger() *slog.Logger {
	return logger
}
