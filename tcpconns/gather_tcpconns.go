package tcpconns

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	TcpNull = iota
	TcpEstablished
	TcpSynSent
	TcpSynRecv
	TcpFinWait1
	TcpFinWait2
	TcpTimeWait
	TcpClose
	TcpCloseWait
	TcpLastAck
	TcpListen
	TcpClosing
	TcpMaxStates
)

var TcpStates = map[int]string{
	TcpEstablished: "ESTABLISHED",
	TcpSynSent:     "SYN_SENT",
	TcpSynRecv:     "SYN_RECV",
	TcpFinWait1:    "FIN_WAIT1",
	TcpFinWait2:    "FIN_WAIT2",
	TcpTimeWait:    "TIME_WAIT",
	TcpClose:       "CLOSED",
	TcpCloseWait:   "CLOSE_WAIT",
	TcpLastAck:     "LAST_ACK",
	TcpListen:      "LISTEN",
	TcpClosing:     "CLOSING",
}

type ConnResult struct {
	Connections map[uint16]*TcpConnection
}

type TcpConnection struct {
	Conns  int32
	Status map[string]int32
}

func GatherTcpconnsInfo(connpath string) (*ConnResult, error) {
	p := new(ConnResult)
	p.Connections = make(map[uint16]*TcpConnection)

	// use ReadAll because otherwise /proc/net/tcp may change from
	// underneath us
	var buf []byte
	for {
		connstat, err := os.Open(connpath)
		if err != nil {
			return nil, err
		}
		defer connstat.Close()
		buf, err = ioutil.ReadAll(connstat)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				continue
			}
			return nil, err
		}
		break
	}
	lines := bytes.Split(buf, []byte("\n"))
	for _, l := range lines {
		l = bytes.Trim(l, "\x00")
		c := strings.Fields(string(l))
		if len(c) == 0 || c[0] == "sl" {
			continue
		}
		port, err := getPort(c[1])
		if err != nil {
			return nil, err
		}
		status, err := strconv.ParseInt(c[3], 16, 32)
		if err != nil {
			return nil, err
		}
		if _, ok := p.Connections[port]; !ok {
			p.Connections[port] = new(TcpConnection)
			p.Connections[port].Status = make(map[string]int32)
		}
		p.Connections[port].Conns++
		p.Connections[port].Status[TcpStates[int(status)]]++
	}

	return p, nil
}

func getPort(addr string) (uint16, error) {
	re := regexp.MustCompile(`[0-9A-F]+:([0-9A-F]+)`)
	p := re.FindStringSubmatch(addr)
	if p == nil {
		err := fmt.Errorf("no port found in '%s'", addr)
		return 0, err
	}
	port, err := strconv.ParseUint(p[1], 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}
