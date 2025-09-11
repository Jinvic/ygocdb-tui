package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"ygocdb-tui/internal/log"
	"ygocdb-tui/internal/ui"
)

// logLevel represents the log level flag
type logLevel struct {
	value log.LogLevel
	set   bool
}

// String returns the string representation of the log level
func (l *logLevel) String() string {
	return l.value.String()
}

// Set sets the log level from a string
func (l *logLevel) Set(value string) error {
	switch value {
	case "off", "OFF":
		l.value = log.OffLevel
	case "error", "ERROR":
		l.value = log.ErrorLevel
	case "warn", "WARN":
		l.value = log.WarnLevel
	case "info", "INFO":
		l.value = log.InfoLevel
	case "debug", "DEBUG":
		l.value = log.DebugLevel
	default:
		return fmt.Errorf("invalid log level: %s", value)
	}
	l.set = true
	return nil
}

func main() {
	// Define command line flags
	var logLevelFlag logLevel
	flag.Var(&logLevelFlag, "log-level", "set log level (off, error, warn, info, debug)")
	
	// Set usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	
	// Parse command line flags
	flag.Parse()
	
	// Initialize logger if log level is set
	if logLevelFlag.set {
		if err := log.Init(logLevelFlag.value); err != nil {
			stdlog.Printf("failed to initialize logger: %v", err)
			os.Exit(1)
		}
		defer log.Close()
		
		// Log application start
		log.Info("ygocdb-tui started with log level: %s", logLevelFlag.value.String())
	}
	
	// Start the TUI application
	if err := ui.Start(); err != nil {
		if logLevelFlag.set {
			log.Error("application error: %v", err)
			log.Close()
		}
		stdlog.Fatal(err)
	}
	
	// Log application exit
	if logLevelFlag.set {
		log.Info("ygocdb-tui exited normally")
	}
}