package cryptogen

import (
	"fmt"
	"os/exec"
)

func GenerateCryptoConfig() {
	cmd := exec.Command("cryptogen", "generate", "--config=../../outputs/crypto-config.yaml")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return
	}
	fmt.Println("Execute Command finished.")
}
