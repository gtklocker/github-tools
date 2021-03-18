package main

import (
	"log"
	"os"
	"context"
	"flag"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var orgToNuke = flag.String("orgToNuke", "", "the name of the organization you want to stop receiving notifications from")

func main() {
	flag.Parse()
	if len(*orgToNuke) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Printf("we will get rid of notifications from %s", *orgToNuke)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GHACCESSTOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	for pg := 1; ; pg++ {
		repos, _, err := client.Activity.ListWatched(ctx, "", &github.ListOptions{Page: pg, PerPage: 100})

		if err == nil {
			if len(repos) == 0 {
				break
			}
			for _, r := range repos {
				owner := r.GetOwner().GetLogin()
				repoName := r.GetName()
				if owner == *orgToNuke {
					log.Printf("unwatching %s/%s", owner, repoName)
					_, err := client.Activity.DeleteRepositorySubscription(ctx, owner, repoName)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			log.Fatal(err)
		}
	}
	log.Print("done!")
}
