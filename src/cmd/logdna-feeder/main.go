package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ctrlrsf/logdna"
	// TODO: promote this to $GOPATH?
	"../../internal/feeder"
)

// this is a comment

func main() {

	apiKey := os.Getenv("LOGDNA_API_KEY")

	if apiKey == "" {
		fmt.Println("Set LOGDNA_API_KEY env var")
		os.Exit(1)
	}

	hostname := flag.String("hostname", "", "hostname you want logs to appear from in LogDNA viewer")
	logFileName := flag.String("log-file-name", "", "log file or app name you want logs to appear as in LogDNA viewer")

	flag.Parse()

	if *hostname == "" {
		fmt.Println("Error: hostname flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if *logFileName == "" {
		fmt.Println("Error: log-file-name flag is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg := logdna.Config{
		APIKey:     apiKey,
		LogFile:    *logFileName,
		Hostname:   *hostname,
		FlushLimit: logdna.DefaultFlushLimit,
	}
	client := logdna.NewClient(cfg)

	logdnafeeder := feeder.NewLogdnaFeeder( feeder.NewTemplateGenerator(), client)

	fmt.Printf("Feed pid=%d\n", os.Getpid())

	tick, _ := time.ParseDuration("500ms")
	logsPerTick := 1
	nLogsToSend := 10
	logdnafeeder.SendLogs(tick, logsPerTick, nLogsToSend)

}
