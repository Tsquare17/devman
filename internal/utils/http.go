package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, response.Body)

	return err
}
