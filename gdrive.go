package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type WriteCounter struct {
	Total    uint64
	Filename string
}

func GetConfirmCodeAndCookies(exportUrl string) ([]*http.Cookie, string, error) {
	resp, err := http.Get(exportUrl)
	if err != nil {
		return nil, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, "", err
	}

	ucHtmlLink := htmlquery.FindOne(doc, "//a[@id='uc-download-link']")
	ucLink := htmlquery.SelectAttr(ucHtmlLink, "href")

	ucConfirmCode := strings.Split(strings.Split(ucLink, "&")[1], "=")[1]

	return resp.Cookies(), ucConfirmCode, nil
}

func GdriveDownload(url, filePath string) error {
	// get file id
	splitted := strings.Split(url, "/")
	fileId := splitted[5]

	exportUrl := fmt.Sprintf("https://drive.google.com/uc?id=%v&export=download", fileId)

	cookies, confirmCode, err := GetConfirmCodeAndCookies(exportUrl)
	if err != nil {
		return err
	}

	confirmUrl := fmt.Sprintf("%v&confirm=%v", exportUrl, confirmCode)

	err = DownloadFile(filePath, confirmUrl, cookies)

	return err
}

func DownloadFile(filepath string, url string, cookies []*http.Cookie) error {
	client := http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for _, c := range cookies {
		req.AddCookie(c)
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

	// Write the body to file
	counter := &WriteCounter{Filename: filepath}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))

	fmt.Println()

	return err
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading %s... %s complete", wc.Filename, humanize.Bytes(wc.Total))
}
