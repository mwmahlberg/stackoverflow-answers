package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"log/syslog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

var (
	interval     time.Duration = 2 * time.Second
	host         string        = "localhost:514"
	network      string        = "tcp"
	networkRegex               = regexp.MustCompile("^(tcp|udp)$")
)

func init() {
	flag.DurationVar(&interval, "interval", interval, "Interval between log messages")
	flag.StringVar(&host, "host", host, "Host to send logs to")
	flag.StringVar(&network, "network", network, "Network protocol to use (tcp/udp)")
}

func main() {
	flag.Parse()

	// Validate network protocol
	if !networkRegex.MatchString(network) {
		panic(fmt.Errorf("invalid network protocol: %s, must be tcp or udp", network))
	}

	w, err := syslog.Dial(network, host, syslog.LOG_INFO|syslog.LOG_LOCAL0, "test")
	// var plogger psyslog.Logger
	if err != nil {
		panic(fmt.Errorf("failed to connect to syslog server: %w", err))
	}
	defer w.Close()
	slog.Info("Connected to syslog server", "host", host, "network", network)

	// We want to shut down gracefully
	// so we need to handle signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start logging messages at the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	ourname, err := os.Hostname()
	if err != nil {
		slog.Error("Failed to get hostname", "error", err)
		ourname = "unknown"
	}
	for {
		select {
		case <-ctx.Done():
			if err := w.Warning(fmt.Sprintf("%s %s %s", ourname, os.Args[0], "Shutting down gracefully")); err != nil {
				slog.Error("Failed to write shutdown message", "error", err)
			} else {
				slog.Info("Shutdown message sent successfully")
			}
			slog.Info("Shut down gracefully")
			return
		case <-ticker.C:
			slog.Info("Logging message at interval", "interval", interval)

			if err := w.Info(fmt.Sprintf("%s %s %s", ourname, os.Args[0], "Hello, World!")); err != nil {
				slog.Error("Failed to write message", "error", err)
			} else {
				slog.Info("Message sent successfully")
			}
		}
	}
}
