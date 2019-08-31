package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"strings"
)

type Episode struct {
	Title      string
	EpisodeNr  string
	GDriveLink string
}

func DownloadEpisode(episode *Episode) error {
	// get current working directory
	wcdir, err := os.Getwd()
	if err != nil {
		return err
	}

	// replace spaces in episode title with dots (.)
	episodeTitle := strings.ReplaceAll(episode.Title, " ", ".")

	episodeFileName := fmt.Sprintf("%v-%v.mp4", episode.EpisodeNr, episodeTitle)
	rarDownloadPath := fmt.Sprintf("/tmp/%v.rar", episode.EpisodeNr)
	extractionPath := fmt.Sprintf("/tmp/%v-extracted", episode.EpisodeNr)

	// Download rar archived video into /tmp directory
	err = GdriveDownload(episode.GDriveLink, rarDownloadPath)
	if err != nil {
		return err
	}

	// extract rar archive in /tmp directory
	fmt.Printf("Extracting %v to %v\n", rarDownloadPath, extractionPath)
	err = archiver.Unarchive(rarDownloadPath, extractionPath)

	// get video filename and rename it
	files, err := ioutil.ReadDir(extractionPath)
	videoName := files[0].Name()
	err = os.Rename(fmt.Sprintf("%v/%v", extractionPath, videoName), fmt.Sprintf("%v/%v", wcdir, episodeFileName))
	if err != nil {
		return err
	}
	fmt.Println("Copy video to your path")

	// remove rar file and extracted directory
	fmt.Println("Removing temporary files")
	err = os.Remove(rarDownloadPath)
	err = os.RemoveAll(extractionPath)

	return err
}

func FetchEpisodes() ([]*Episode, error) {
	url := "https://amalgam-fansubs.moe/detektiv-conan-2017/"
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	conanDiv := htmlquery.FindOne(doc, "//div[@id='conan']")
	episodeTable := htmlquery.FindOne(conanDiv, "table")

	rows := htmlquery.Find(episodeTable, "//tr")

	var episodes []*Episode
	for i := 1; i < len(rows); i++ {
		cols := htmlquery.Find(rows[i], "//td")
		episodeNr := htmlquery.InnerText(cols[0]) // number
		episodeNr = strings.TrimSuffix(episodeNr, ".")
		episodeTitle := htmlquery.InnerText(cols[1])                   // title
		gdriveLink := htmlquery.SelectAttr(cols[3].FirstChild, "href") // gdrive link

		episodes = append(episodes, &Episode{
			Title:      episodeTitle,
			EpisodeNr:  episodeNr,
			GDriveLink: gdriveLink,
		})
	}

	return episodes, nil
}
