package main

import (
	"amalgamDCLoader/amalgam"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
)

func CmdApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Detektiv Conan Amalgam Downloader"
	app.Usage = "Download Detektiv Conan Episodes from https://amalgam-fansubs.moe/"

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list available episodes",
			Action:  ListEpisodes,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:     "dlink,d",
					Usage:    "List episodes with download links",
					Required: false,
					Hidden:   false,
				},
				cli.BoolFlag{
					Name:     "gdrive,g",
					Usage:    "Show if episodes can be downloaded via google drive",
					Required: false,
					Hidden:   false,
				},
				cli.StringFlag{
					Name:     "format",
					Usage:    "available values: csv, html, md",
					Required: false,
					Hidden:   false,
					Value:    "",
				},
			},
		},
		{
			Name:      "download",
			Aliases:   []string{"d"},
			Usage:     "download episodes",
			ArgsUsage: "episode list: 1 2 3  episode range: 4-10, combined: 1 2-6 8",
			Action:    DownloadEpisodes,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:     "gdrive,g",
					Usage:    "Download episode from google drive",
					Required: false,
					Hidden:   false,
				},
			},
		},
	}

	return app
}

func ListEpisodes(c *cli.Context) error {
	episodes, err := amalgam.FetchEpisodes()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	if c.Bool("dlink") {
		t.AppendHeader(table.Row{"Nr.", "Title", "Download Link", "Google Drive Link"})
	} else if c.Bool("gdrive") {
		t.AppendHeader(table.Row{"Nr.", "Title", "Google Drive"})
	} else {
		t.AppendHeader(table.Row{"Nr.", "Title"})
	}

	for _, e := range episodes {
		gdAvailable := "✓"
		if e.GDriveLink == "" {
			gdAvailable = "✘"
		}

		if strings.Contains(e.EpisodeNr, ",") {
			e.Title = fmt.Sprintf("%v (Combined Episode)", e.Title)
		}

		// Skip if available flag and episode is not downloadable
		if c.Bool("gdrive") && e.GDriveLink == "" {
			continue
		}

		if c.Bool("dlink") {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title, e.DownloadLink, e.GDriveLink})
		} else if c.Bool("gdrive") {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title, gdAvailable})
		} else {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title})
		}

	}

	t.AppendFooter(table.Row{fmt.Sprintf("Total: %v", len(episodes))})

	if c.String("format") == "csv" {
		t.RenderCSV()
	} else if c.String("format") == "html" {
		t.RenderHTML()
	} else if c.String("format") == "md" {
		t.RenderMarkdown()
	} else {
		t.Render()
	}

	return nil
}

func DownloadEpisodes(c *cli.Context) error {
	var episodeArgList []string
	episodesArgs := c.Args()

	// Fetch episodes and create map for easier download
	episodesList, err := amalgam.FetchEpisodes()
	if err != nil {
		return err
	}
	episodes := make(map[string]*amalgam.Episode)
	for _, e := range episodesList {
		episodes[e.EpisodeNr] = e
	}

	// Parse episodes argument
	for _, s := range episodesArgs {
		if strings.Contains(s, "-") {
			// handle ranges
			splitted := strings.Split(s, "-")
			start, err := strconv.Atoi(splitted[0])
			if err != nil {
				return err
			}
			end, err := strconv.Atoi(splitted[1])
			if err != nil {
				return err
			}
			for i := start; i <= end; i++ {
				episodeArgList = append(episodeArgList, strconv.Itoa(i))
				// handle ,5 episodes (there are episodes with numbers like 704,5)
				//if _, ok := episodes[fmt.Sprintf("%v,5", i)]; ok {
				//	episodeArgList = append(episodeArgList, fmt.Sprintf("%v,5", i))
				//}
			}
		} else {
			// handle simple comma separation
			episodeArgList = append(episodeArgList, s)
		}
	}

	fmt.Println("Downloading Episodes:", strings.Join(episodeArgList, " "))
	fmt.Println()

	// download episodes
	for _, episodeNr := range episodeArgList {
		// check if episode is available
		episode, ok := episodes[episodeNr]
		if !ok {
			fmt.Printf("Episode %v is not available\n", episodeNr)
			continue
		}

		if c.Bool("gdrive") && episode.GDriveLink == "" {
			fmt.Printf("Episode %v is not available via Google Drive\n", episodeNr)
			continue
		}

		if c.Bool("gdrive") {
			err = amalgam.DownloadEpisodeFromGDrive(episode)
		} else {
			err = amalgam.DownloadEpisode(episode)
		}

		if err != nil {
			fmt.Printf("Error while downloading Episode %v\n", episodeNr)
		}

		fmt.Println()
	}

	return nil
}
