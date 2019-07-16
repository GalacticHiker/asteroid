package feeder

import (
	"time"
)

// Feed data and
type Feed struct {
	lg LogGenerator
}

// Feeder interface
type Feeder interface {
	SendLogs(tick time.Duration, logsPerTick, nLogsToSend int)
	logGenerator() LogGenerator
	writeLog(logText string) (bytesSent int64)
}

type feedState struct {
	nLogsToSend   int
	eof           bool
	logsPerTick int

	logsSentCount int
	logBytesSentCount int64

	clockStart   time.Time
	tickStart    time.Time
	tickDuration time.Duration

}

func newFeedState(nLogsToSend, logsPerTick int) *feedState {

	fs := new(feedState)
	fs.nLogsToSend = nLogsToSend
	fs.logsPerTick = logsPerTick

	fs.clockStart = time.Now()

	return fs
}
func sendClock(f Feeder, tick time.Duration, logsPerTick, nLogsToSend int) {

	fs := newFeedState(nLogsToSend, logsPerTick)

	defer func() {
		time.Sleep(3 * time.Second) // drain buffers
	}()

	for {

		for {

			fs.tick(f)

			// do not wait for a tick when it takes longer than a tick to send
			if fs.tickDuration > tick {
				continue
			}
			break
		}

		if fs.allLogsSent() {
			break
		}

		select {
		case <-time.Tick(tick):
		}

	}

}

func (fs *feedState) tick(f Feeder) {

	nLogsToSendOnTick := fs.nLogsToSend - fs.logsSentCount
	if nLogsToSendOnTick >= fs.logsPerTick {
		nLogsToSendOnTick = fs.logsPerTick
	}

	fs.tickStart = time.Now()

	logsSent := 0
	bytesSent := int64(0)
	eof := false
	if nLogsToSendOnTick > 0 {
		bytesSent, logsSent, eof = sendNLogs(f, nLogsToSendOnTick)
	}

	fs.logsSentCount += logsSent
	fs.logBytesSentCount += bytesSent
	fs.eof = eof

	return
}

// sendNLogs sends the specified number of logs
func sendNLogs(f Feeder, nLogsToSendOnTick int) (int64, int, bool) {

	eof := false 

	bytesSent := int64(0)
	logsSent := 0

	for ; logsSent < nLogsToSendOnTick; logsSent++ {

		log, _ := f.logGenerator().GenerateLog(time.Now())
		if eof {
			break
		}

		bytesSent = f.writeLog(log)

	}

	return bytesSent, logsSent, eof
}
func (fs *feedState) allLogsSent() bool {
	if fs.logsSentCount == fs.nLogsToSend || fs.eof {
		return true
	}

	return false
}
