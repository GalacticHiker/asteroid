package feeder

import (
	"log"
	"net"
	"time"
)

// TCPSyslogConf configuration 
type TCPSyslogConf struct {
	DestAddr string
}

// NewTCPSyslogConf return the default TCP configuration
func NewTCPSyslogConf() *TCPSyslogConf {

	return &TCPSyslogConf {
		DestAddr :  "",
	}

}

// TCPFeeder struct
type TCPFeeder struct {
	Feed
	tcpconn net.Conn
}

// NewTCPFeeder create a TCPfeeder
func NewTCPFeeder(lg LogGenerator, conf *TCPSyslogConf) *TCPFeeder {

	f := new(TCPFeeder)

	f.lg = lg

	c, err := net.Dial("tcp", conf.DestAddr)
	f.tcpconn = c
	if err != nil {
		log.Fatalf("Error Connecting. %v. Is rsyslogd listening?\n", err)
	}

	return f
}

func (f *TCPFeeder) logGenerator() LogGenerator {
	return f.lg
}

func (f *TCPFeeder) close() {
	f.tcpconn.Close()
}

// SendLogs ...
func (f *TCPFeeder) SendLogs(tick time.Duration, logsPerTick int, nLogsToSend int) {

	sendClock(f, tick, logsPerTick, nLogsToSend)
	f.close()

}

func (f *TCPFeeder) writeLog(logText string) (bytesSent int64) {
	payload := []byte(logText)

	n, err := f.tcpconn.Write(payload)

	if err != nil {
		log.Fatalf("Error sending: %v\n", err)
	}

	if n != len(payload) {
		log.Fatalf("PayloadLength==%d only %d sent.", len(payload), n)
	}

	return int64(n)
}
