package configtx

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type AnchorPeer struct {
	Host string
	Port int
}

type Organization struct {
	Name       string
	AnchorPeer AnchorPeer
}

func (o Organization) GetNameToLower() string {
	return strings.ToLower(o.Name)
}

type BatchSize struct {
	MaxMessageCount   int
	AbsoluteMaxBytes  int
	PreferredMaxBytes int
}

type Orderer struct {
	Host string
	Port int
}

type Consensus struct {
	OrdererType  string
	BatchTimeout int32
	BatchSize    BatchSize
	Orderers     []Orderer
}

type Channel struct {
	Name           string
	Consortium     string
	Organizatioins []Organization
}

type Configtx struct {
	Organizations []Organization
	Consensus     Consensus
	Channel       Channel
}

//GenerateConfigTxTemplate To write the config to the w according to the template tpl
func GenerateConfigTxTemplate(configtx Configtx, tpl *template.Template, w io.Writer) error {
	err := tpl.Execute(w, configtx)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
