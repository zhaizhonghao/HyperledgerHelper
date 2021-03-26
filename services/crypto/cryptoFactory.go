package crypto

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
)

type OrdererCp struct {
	HostName string `json:"HostName"`
}

type PeerOrgCp struct {
	Name         string `json:"Name"`
	Domain       string `json:"Domain"`
	CountOfPeers int32  `json:"CountOfPeers"`
	CountOfUsers int32  `json:"CountOfUsers"`
}

func (p PeerOrgCp) GenListFromPeers() []int {
	ret := make([]int, p.CountOfPeers)
	for i := 0; i < (int)(p.CountOfPeers); i++ {
		ret[i] = i
	}
	return ret
}

func (p PeerOrgCp) GetIdOfPeerOrg() int {
	var part = strings.Split(p.Name, "g")
	port, err := strconv.Atoi(part[1])
	if err != nil {
		fmt.Println(err)
	}
	return port
}

func (p PeerOrgCp) GetNameToLower() string {

	return strings.ToLower(p.Name)
}

func (p PeerOrgCp) GetPortOfPeer() int {
	var part = strings.Split(p.Name, "g")
	var port int
	if part[1] == "" {
		port = 1
	} else {
		portTemp, err := strconv.Atoi(part[1])
		if err != nil {
			fmt.Println(err)
			return -1
		}
		port = portTemp

	}
	return (port+6)*1000 + 51
}

func (p PeerOrgCp) GetPortOfCA() int {
	var part = strings.Split(p.Name, "g")
	port, err := strconv.Atoi(part[1])
	if err != nil {
		fmt.Println(err)
	}
	return (port+6)*1000 + 54
}

func (o OrdererCp) GetGeneralPortOfOrderer() int {
	var part = strings.Split(o.HostName, "erer")
	var port int
	if part[1] == "" {
		port = 1
	} else {
		portTemp, err := strconv.Atoi(part[1])
		if err != nil {
			fmt.Println(err)
			return -1
		}
		port = portTemp

	}

	return (port+6)*1000 + 50
}

func (o OrdererCp) GetOperationPortOfOrderer() int {
	var part = strings.Split(o.HostName, "erer")
	var port int
	if part[1] == "" {
		port = 1
	} else {
		portTemp, err := strconv.Atoi(part[1])
		if err != nil {
			fmt.Println(err)
			return -1
		}
		port = portTemp

	}
	return 8442 + port
}

type ConfigCp struct {
	OrdererCps []OrdererCp `json:"OrdererCps"`
	PeerOrgCps []PeerOrgCp `json:"PeerOrgCps"`
}

//GenerateConfigTxTemplate To write the config to the w according to the template tpl
func GenerateCryptoTemplate(configCp ConfigCp, tpl *template.Template, w io.Writer) error {
	err := tpl.Execute(w, configCp)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
