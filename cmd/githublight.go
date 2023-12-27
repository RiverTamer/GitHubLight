package main

import (
	"context"
	"github.com/google/go-github/v57/github"
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITHUBLIGHT_ACCESS_TOKEN")
	if len(token) == 0 {
		log.Fatalf("You need to define GITHUBLIGHT_ACCESS_TOKEN.")
	}
	client := github.NewClient(nil).WithAuthToken(token)
	options := github.SearchOptions{
		Sort:      "committer-date",
		Order:     "desc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}
	commits, _, err := client.Search.Commits(context.Background(), "author:KarlKraft", &options)
	if err != nil {
		log.Fatalf("Could not fetch commits. %v", err)
		return
	}
	log.Printf("Total commit count is %d", *commits.Total)

}
