package logmill

import (
	"time"

	"github.com/ctrlrsf/logdna"
)

// LogdnaConf configuration 
type LogdnaConf struct {
	Hostname string
	LogFilename string
}

// NewLogdnaConf return the default logdna configuration
func NewLogdnaConf() *LogdnaConf {

	return &LogdnaConf {
		Hostname :  "logmill-logdna",
		LogFilename : "logmill-file",
	}

}

// LogdnaLogmill sends logs to the logdna server
type LogdnaLogmill struct {
	Mill
	client *logdna.Client
}

// NewLogdnaLogmill creates a new Logdna Logmill
func NewLogdnaLogmill( lg LogGenerator, client *logdna.Client) *LogdnaLogmill {

	logdnaLogmill := new(LogdnaLogmill)

	logdnaLogmill.lg = lg

	logdnaLogmill.client = client
	return logdnaLogmill
}

// SendLogs send the logs per the configuration
func (f *LogdnaLogmill) SendLogs( tick time.Duration, logsPerTick int , nLogsToSend int ) {


	sendClock(f, tick, logsPerTick, nLogsToSend)

	f.client.Flush() // TODO: necc? Close() flushes?
	f.client.Close()

}

func (f *LogdnaLogmill) logGenerator() LogGenerator {
	return f.lg
}

func (f *LogdnaLogmill) writeLog(logText string) (bytesSent int64) {

	t := f.lg.TemplateContext().LogTime
	f.client.Log(t, logText)

	return int64(len(logText))
}
