package scan

import (
	"fmt"
	"net"
	"time"
)

// represents state of single TCP port
type PortState struct {
	Port int
	Open state
}

type state bool

// represents scan results for a single host
type Results struct {
	Host       string
	NotFound   bool        //indicates whether host can be resolved to a valid IP addrss or not
	PortStates []PortState //indicates status of each port scanned
}

// converts boolean value of state to a reable string
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// other packages cant use this function directly--ITS PRIVATE!!!
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}
	// JoinHostPort combines host and port into a network address of the form "host:port"
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	// dialtimeout tries to connect to a address within a given time
	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true

	return p
}

// performs a port scan on the hosts list
func Run(hl *HostLists, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}
	
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}