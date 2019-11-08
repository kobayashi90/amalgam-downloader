package lib

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type WriteCounter struct {
	Current  uint64
	Total    uint64
	Filename string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += uint64(n)

	wc.PrintProgress()

	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	if wc.Total == 0 {
		fmt.Printf("\rDownloading %s: %s", wc.Filename, humanize.Bytes(wc.Current))
	} else {
		fmt.Printf("\rDownloading %s: %s / %s (%v %%)", wc.Filename, humanize.Bytes(wc.Current), humanize.Bytes(wc.Total), wc.Current*100/wc.Total)
	}
}

func DownloadFile(filepath string, url string) error {
	client := http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	totalFileSize := resp.ContentLength

	// Write the body to file
	counter := &WriteCounter{Filename: path.Base(filepath), Total: uint64(totalFileSize)}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	return err
}
