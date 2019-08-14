package logmill

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

// TCPLogmill struct
type TCPLogmill struct {
	Mill
	tcpconn net.Conn
}

// NewTCPLogmill create a TCPLogmill
func NewTCPLogmill(lg LogGenerator, conf *TCPSyslogConf) *TCPLogmill {

	f := new(TCPLogmill)

	f.lg = lg

	c, err := net.Dial("tcp", conf.DestAddr)
	f.tcpconn = c
	if err != nil {
		log.Fatalf("Error Connecting. %v. Is rsyslogd listening?\n", err)
	}

	return f
}

func (f *TCPLogmill) logGenerator() LogGenerator {
	return f.lg
}

func (f *TCPLogmill) close() {
	f.tcpconn.Close()
}

// SendLogs ...
func (f *TCPLogmill) SendLogs(tick time.Duration, logsPerTick int, nLogsToSend int) {

	sendClock(f, tick, logsPerTick, nLogsToSend)
	f.close()

}

func (f *TCPLogmill) writeLog(logText string) (bytesSent int64) {
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
