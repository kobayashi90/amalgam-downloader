package main

import (
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
		},
		{
			Name:      "download",
			Aliases:   []string{"d"},
			Usage:     "download episodes",
			ArgsUsage: "episode list: 1,2,3, episode range: 4-10, combined: 1,2-6,8",
			Action:    DownloadEpisodes,
		},
	}

	return app
}

func ListEpisodes(c *cli.Context) error {
	episodes, err := FetchEpisodes()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"Nr.", "Title", "Available"})
	notAvailable := 0
	for _, e := range episodes {
		available := "✓"
		if e.GDriveLink == "" {
			available = "✘"
			notAvailable++
		}
		t.AppendRow(table.Row{e.EpisodeNr, e.Title, available})
	}

	t.AppendFooter(table.Row{fmt.Sprintf("Total: %v", len(episodes)), "", fmt.Sprintf("Available: %v", len(episodes)-notAvailable)})
	t.Render()

	return nil
}

func DownloadEpisodes(c *cli.Context) error {
	var episodeArgList []string
	episodesArg := c.Args().Get(0)

	// Parse episodes argument
	commaSplitted := strings.Split(episodesArg, ",")
	for _, s := range commaSplitted {
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
			}
		} else {
			// handle simple comma separation
			episodeArgList = append(episodeArgList, s)
		}
	}

	// Fetch episodes and create map for easier download
	episodesList, err := FetchEpisodes()
	if err != nil {
		return err
	}
	episodes := make(map[string]*Episode)
	for _, e := range episodesList {
		episodes[e.EpisodeNr] = e
	}

	// download episodes
	for _, episodeNr := range episodeArgList {
		// check if episode is available
		_, ok := episodes[episodeNr]
		if !ok || episodes[episodeNr].GDriveLink == "" {
			fmt.Printf("Episode %v is not available\n", episodeNr)
			continue
		}
		err = DownloadEpisode(episodes[episodeNr])
		if err != nil {
			fmt.Printf("Error while downloading Episode %v\n", episodeNr)
		}

		fmt.Println()
	}

	return nil
}
