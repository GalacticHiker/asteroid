package feeder

import (
	"time"

	"github.com/ctrlrsf/logdna"
)

// LogdnaFeeder feeds logs to the logdna server
type LogdnaFeeder struct {
	Feed
	client *logdna.Client
}

// NewLogdnaFeeder creates a new Logdna feeder
func NewLogdnaFeeder( lg LogGenerator, client *logdna.Client) *LogdnaFeeder {

	logdnafeeder := new(LogdnaFeeder)

	logdnafeeder.lg = lg

	logdnafeeder.client = client
	return logdnafeeder
}

// SendLogs send the logs per the configuration
func (f *LogdnaFeeder) SendLogs( tick time.Duration, logsPerTick int , nLogsToSend int ) {


	sendClock(f, tick, logsPerTick, nLogsToSend)

	f.client.Flush() // TODO: necc? Close() flushes?
	f.client.Close()

}

func (f *LogdnaFeeder) logGenerator() LogGenerator {
	return f.lg
}

func (f *LogdnaFeeder) writeLog(logText string) (bytesSent int64) {

	t := f.lg.TemplateContext().LogTime
	f.client.Log(t, logText)

	return int64(len(logText))
}
