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
	"../../internal/logmill"
)

// this is a comment
func main() {

	// protocol 
	protocol := flag.String("protocol", "" , "tcp | logdna | kafka ")

	// tcp
	tcpSyslogConf := logmill.NewTCPSyslogConf()
	tcpSyslogConf.DestAddr = "localhost:514" 		//default
	flag.CommandLine.StringVar(&tcpSyslogConf.DestAddr, "destAddr", tcpSyslogConf.DestAddr, "Destination IP:port of rsyslog server")

	// logdna
	logdnaConf := logmill.NewLogdnaConf()
	flag.CommandLine.StringVar(&logdnaConf.Hostname, "hostname", logdnaConf.Hostname, "hostname you want logs to appear from in LogDNA viewer")
	flag.CommandLine.StringVar(&logdnaConf.LogFilename, "logdna-file", logdnaConf.LogFilename, "log file or app name you want logs to appear as in LogDNA viewer")

	// kafka
	kafkaConf := logmill.NewKafkaConf();
	flag.CommandLine.StringVar(&kafkaConf.Broker, "broker", kafkaConf.Broker, "Kafka Broker")
	flag.CommandLine.StringVar(&kafkaConf.Topic, "topic", kafkaConf.Topic, "Kafka Topic")

	// rate TODO: create a struct for these values
	tick := flag.Duration("tick", time.Duration(1)*time.Second, "Send frequency")
	logsPerTick := flag.Int("logsPerTick", 1, "Number of logs to send per tick.")
	nLogsToSend := flag.Int("nLogsToSend", 1, "Number of logs to send.")

	// format template
	logTemplate := flag.String("template", "defaultKVP", "Name of logTemplate")

	// tag 
	tag := flag.String("tag", "" , "all logs have this value")

	flag.Parse()

	if *protocol == "" {
		log.Fatalf("Protocol must be tcp | logdna | kafka\n")
	}

	setExeHome()

	fmt.Printf("Logmill pid=%d tag=%s\n", os.Getpid(), *tag)

	// TODO: distinguish correct mill type based on arguments -- create a mill factory
	lg := logmill.NewTemplateGenerator(*logTemplate)
	tc := lg.TemplateContext();
	tc.Tag = *tag
	
	var lm logmill.Logmill
	if *protocol == "tcp" {
		tc.Protocol = "tcp"
		lm = logmill.NewTCPLogmill(lg, tcpSyslogConf)
	} else if *protocol == "logdna" {
		tc.Protocol = "logdna"
		lm = createLogdnaMill( lg, logdnaConf)
	} else if *protocol == "kafka" {
		tc.Protocol = "kafka"
		lm = logmill.NewKafkaLogmill(kafkaConf,lg)
	}

	lm.SendLogs(*tick, *logsPerTick, *nLogsToSend)

}

func setExeHome() {
	// TODO : put this somewhere more appropriate
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("ERROR:%v\n", err)
	}
	bindir := path.Dir(ex) // trim executable name
	os.Chdir(bindir)       // set exe home	
}

func createLogdnaMill( lg logmill.LogGenerator, conf *logmill.LogdnaConf) logmill.Logmill {

	apiKey := os.Getenv("LOGDNA_API_KEY")

	if apiKey == "" {
		fmt.Println("Set LOGDNA_API_KEY env var")
		os.Exit(1)
	}
	if conf.Hostname == "" {
		fmt.Println("Error: hostname flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if conf.LogFilename == "" {
		fmt.Println("Error: log-file-name flag is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg := logdna.Config{
		APIKey:     apiKey,
		LogFile:    conf.LogFilename,
		Hostname:   conf.Hostname,
		FlushLimit: logdna.DefaultFlushLimit,
	}
	client := logdna.NewClient(cfg)

	logdnaLogmill := logmill.NewLogdnaLogmill(lg, client)

	return logdnaLogmill
}