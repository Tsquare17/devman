package utils

import (
	"log"
	"os"
	"os/exec"
	"regexp"
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

func Slugify(title string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Panic(err)
	}

	return reg.ReplaceAllString(title, "_")
}
