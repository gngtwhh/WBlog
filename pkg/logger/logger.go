package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Options struct {
	Level     slog.Level
	FilePath  string // logger file(nil will print to stderr)
	AddSource bool   // record line number
}

func Setup(opts *Options) *slog.Logger {
	var w io.Writer = os.Stderr

	if opts.FilePath != "" {
		dir := filepath.Dir(opts.FilePath)
		_ = os.MkdirAll(dir, 0755)

		file, err := os.OpenFile(opts.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			w = io.MultiWriter(os.Stdout, file)
		}
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, t.Format(time.DateTime))
				}
			}
			return a
		},
	}
	var handler slog.Handler
	if opts.Level == slog.LevelDebug {
		handler = slog.NewTextHandler(w, handlerOpts)
	} else {
		handler = slog.NewJSONHandler(w, handlerOpts)
	}

	logger := slog.New(handler)
	// Set as Global Logger
	slog.SetDefault(logger)

	return logger
}
