package output

import (
	"encoding/xml"
	"fmt"
)

type NmapRun struct {
	Command  string   `xml:"args,attr" json:"command"`
	Start    string   `xml:"start,attr" json:"start"`
	Hosts    []Host   `xml:"host" json:"host"`
	Finished Finished `xml:"runstats>finished" json:"run-stats"`
}

type Host struct {
	Address Address `xml:"address" json:"address"`
	Ports   []Port  `xml:"ports>port" json:"ports"`
}

type Address struct {
	Addr string `xml:"addr,attr" json:"addr"`
	Type string `xml:"addrtype,attr" json:"type"`
}

type Port struct {
	PortId   string `xml:"portid,attr" json:"port-id"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
	State    State  `xml:"state" json:"state"`
}

type State struct {
	State string `xml:"state,attr" json:"state"`
}

type RunStats struct {
	Finished Finished `xml:"finished" json:"finished"`
}

type Finished struct {
	End     string `xml:"time,attr" json:"end"`
	Summary string `xml:"summary,attr" json:"summary"`
}

func (n *NmapRun) IsSuccessfulScan() bool {
	if n.Hosts == nil {
		return false
	}
	return true
}

func ParseRunOutput(out []byte) (*NmapRun, error) {
	var runOutput NmapRun
	err := xml.Unmarshal(out, &runOutput)
	if err != nil {
		return nil, err
	}
	return &runOutput, nil
}

func (n *NmapRun) FormatNmapRun() string {
	var formatted = ""
	for _, host := range n.Hosts {
		formatted += fmt.Sprintf("Host: %v ", host.Address.Addr)
		for _, port := range host.Ports {
			formatted += fmt.Sprintf("Port: %v -> %v ", port.PortId, port.State.State)
		}
		formatted += " | "
	}
	return formatted
}
