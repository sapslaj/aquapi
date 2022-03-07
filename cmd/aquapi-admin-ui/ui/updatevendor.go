package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	fmt.Printf("Downloaded %s to %s\n", url, filepath)
	return err
}

func main() {
	downloadFile("./vendor/bootstrap.min.css", "https://raw.githubusercontent.com/thomaspark/bootswatch/v5/dist/darkly/bootstrap.min.css")
	downloadFile("./vendor/bootstrap.min.css.map", "https://unpkg.com/bootstrap@5.1.3/dist/css/bootstrap.min.css.map")
	downloadFile("./vendor/react.js", "https://unpkg.com/react@17/umd/react.development.js")
	downloadFile("./vendor/react-dom.js", "https://unpkg.com/react-dom@17/umd/react-dom.development.js")
}
