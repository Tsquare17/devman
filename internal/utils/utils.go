package utils

import (
	"os"
	"os/exec"
)

func GetPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}

func GetCommandOutput(command string, args ...string) string {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		panic(err)
	}

	return string(out)
}
