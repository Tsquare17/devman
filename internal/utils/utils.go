package utils

import (
	"os"
)

func GetPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}
