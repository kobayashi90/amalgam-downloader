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
	Title        string
	EpisodeNr    string
	DownloadLink string
	Downloadable bool
	Note         string
}

func DownloadEpisode(episode *Episode) error {
	// get current working directory
	wcdir, err := os.Getwd()
	if err != nil {
		return err
	}

	// replace spaces in episode title with dots (.)
	episodeTitle := strings.ReplaceAll(episode.Title, " ", ".")

	tmpDir := os.TempDir()
	episodeFileName := fmt.Sprintf("%v-%v.mp4", episode.EpisodeNr, episodeTitle)
	rarDownloadPath := fmt.Sprintf("%v/%v.rar", tmpDir, episode.EpisodeNr)
	extractionPath := fmt.Sprintf("%v/%v-extracted", tmpDir, episode.EpisodeNr)

	// Download rar archived video into /tmp directory
	err = GdriveDownload(episode.DownloadLink, rarDownloadPath)
	if err != nil {
		return err
	}

	// extract rar archive in /tmp directory
	fmt.Printf("  --> Extracting %v to %v\n", rarDownloadPath, extractionPath)
	err = archiver.Unarchive(rarDownloadPath, extractionPath)

	// get video filename and rename it
	files, err := ioutil.ReadDir(extractionPath)
	videoName := files[0].Name()
	err = os.Rename(fmt.Sprintf("%v/%v", extractionPath, videoName), fmt.Sprintf("%v/%v", wcdir, episodeFileName))
	if err != nil {
		return err
	}
	fmt.Println("  --> Copy video to your path")

	// remove rar file and extracted directory
	fmt.Println("  --> Removing temporary files")
	err = os.Remove(rarDownloadPath)
	err = os.RemoveAll(extractionPath)

	return err
}

func FetchEpisodes() ([]*Episode, error) {
	urls := []string{
		"https://amalgam-fansubs.moe/detektiv-conan/",
		"https://amalgam-fansubs.moe/detektiv-conan-2017/",
	}

	var episodes []*Episode
	for _, url := range urls {
		doc, err := htmlquery.LoadURL(url)
		if err != nil {
			return nil, err
		}

		conanDiv := htmlquery.FindOne(doc, "//div[@id='conan']")
		episodeTable := htmlquery.FindOne(conanDiv, "table")

		rows := htmlquery.Find(episodeTable, "//tr")

		for i := 1; i < len(rows); i++ {
			cols := htmlquery.Find(rows[i], "//td")
			episodeNr := htmlquery.InnerText(cols[0]) // number
			episodeNr = strings.ReplaceAll(episodeNr, ".", "")
			episodeTitle := htmlquery.InnerText(cols[1])                     // title
			downloadLink := htmlquery.SelectAttr(cols[3].FirstChild, "href") // gdrive link
			downloadable := strings.Contains(downloadLink, "drive.google")

			episodes = append(episodes, &Episode{
				Title:        episodeTitle,
				EpisodeNr:    strings.TrimSpace(episodeNr),
				DownloadLink: downloadLink,
				Downloadable: downloadable,
			})
		}
	}

	return episodes, nil
}
