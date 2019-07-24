package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
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

	// optional
	tick := flag.Duration("tick", time.Duration(1)*time.Second, "Send frequency")
	logsPerTick := flag.Int("logsPerTick", 1, "Number of logs to send per tick.")
	nLogsToSend := flag.Int("nLogsToSend", 1, "Number of logs to send.")

	logTemplate := flag.String("template", "defaultKVP", "Name of logTemplate")
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

	// TODO o: put this somewhere more appropriate
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("ERROR:%v\n", err)
	}
	bindir := path.Dir(ex) // trim executable name
	os.Chdir(bindir)       // set root

	logdnafeeder := feeder.NewLogdnaFeeder(feeder.NewTemplateGenerator(*logTemplate), client)

	fmt.Printf("Feed pid=%d\n", os.Getpid())

	logdnafeeder.SendLogs(*tick, *logsPerTick, *nLogsToSend)

}
