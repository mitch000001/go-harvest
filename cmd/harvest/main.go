package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
	"github.com/mitch000001/go-harvest/harvest/auth"
)

func main() {
	subdomain := os.Getenv("HARVEST_SUBDOMAIN")
	username := os.Getenv("HARVEST_USERNAME")
	password := os.Getenv("HARVEST_PASSWORD")

	clientProvider := auth.NewBasicAuthClientProvider(&auth.BasicAuthConfig{username, password})

	client, err := harvest.New(subdomain, clientProvider.Client)
	if err != nil {
		fmt.Printf("There was an error creating the client:\n")
		fmt.Printf("%T: %v\n", err, err)
		os.Exit(1)
	}
	var projects []*harvest.Project
	err = client.Projects.All(&projects, nil)
	if err != nil {
		fmt.Printf("There was an error fetching all projects:\n")
		fmt.Printf("%T: %v\n", err, err)
		os.Exit(1)
	}
	timeframe := harvest.Timeframe{
		StartDate: harvest.Date(2014, 01, 01, time.UTC),
		EndDate:   harvest.Date(2014, 02, 07, time.UTC),
	}
	params := url.Values{}
	params.Add("from", timeframe.StartDate.Format("2006-01-02"))
	params.Add("to", timeframe.EndDate.Format("2006-01-02"))
	for _, project := range projects {
		fmt.Printf("Project: %+#v\n", project)
		var dayEntries []*harvest.DayEntry
		err := client.Projects.DayEntries(project).All(&dayEntries, params)
		if err != nil {
			fmt.Printf("There was an error fetching all day entries from project with id %d:\n", project.Id())
			fmt.Printf("%T: %v\n", err, err)
		} else {
			for _, d := range dayEntries {
				fmt.Printf("DayEntry: %+#v\n", d)
			}
		}
	}
}
