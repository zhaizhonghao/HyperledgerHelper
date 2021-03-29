package connection

import (
	"fmt"
	"io"
	"text/template"
)

type Channel struct {
	ChannelName string
}

//GenerateConfigTxTemplate To write the config to the w according to the template tpl
func GenerateConnectionTemplate(channel Channel, tpl *template.Template, w io.Writer) error {
	err := tpl.Execute(w, channel)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
