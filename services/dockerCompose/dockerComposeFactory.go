package dockerCompose

import (
	"fmt"
	"html/template"
	"io"

	"github.com/zhaizhonghao/configtxTool/services/crypto"
)

//GenerateConfigTxTemplate To write the config to the w according to the template tpl
func GenerateDockerComposeTemplate(configCp crypto.ConfigCp, tpl *template.Template, w io.Writer) error {
	err := tpl.Execute(w, configCp)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
