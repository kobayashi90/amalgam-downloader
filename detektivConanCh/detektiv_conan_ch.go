package detektivConanCh

import (
	"amalgamDCLoader/lib"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/mholt/archiver"
	"os"
	"path/filepath"
	"strings"
)

type Music struct {
	Title        string
	Filename     string
	DownloadLink string
}

func FetchMusic() ([]*Music, error) {
	musicUrl := "https://www.detektiv-conan.ch/index.php?page=aplayer/musik.php"
	doc, err := htmlquery.LoadURL(musicUrl)
	if err != nil {
		return nil, err
	}

	musicDivs := htmlquery.Find(doc, "//div[@class='album_content']")

	var musics []*Music
	for _, musicDiv := range musicDivs {
		htmlLink := htmlquery.FindOne(musicDiv, "//a")
		relativeDownloadLink := htmlquery.SelectAttr(htmlLink, "href")
		if !strings.Contains(relativeDownloadLink, ".zip") {
			continue
		}
		downloadLink := fmt.Sprintf("https://www.detektiv-conan.ch%v", relativeDownloadLink)

		filename := filepath.Base(relativeDownloadLink)
		title := strings.TrimSuffix(filename, ".zip")

		musics = append(musics, &Music{
			Title:        title,
			Filename:     filename,
			DownloadLink: downloadLink,
		})
	}

	return musics, nil
}

func DownloadMusic(music *Music, unzip, keepArchive bool) error {
	wcdir, err := os.Getwd()
	if err != nil {
		return err
	}

	fp := fmt.Sprintf("%v/%v", wcdir, music.Filename)

	err = lib.DownloadFile(fp, music.DownloadLink)

	if unzip {
		// create directory for extraction
		archivePath := filepath.Dir(fp)
		extractionDirPath := fmt.Sprintf("%s/%s/", archivePath, strings.TrimSuffix(music.Filename, ".zip"))
		err = os.Mkdir(extractionDirPath, 0775)
		if err != nil {
			return err
		}

		err = archiver.Unarchive(fp, extractionDirPath)
		if err != nil {
			return err
		}

		if !keepArchive {
			err = os.Remove(fp)
			if err != nil {
				return err
			}
		}
	}

	return err
}
