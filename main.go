package main

import (
	"log"
	"os"
	"context"
	"flag"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var org = flag.String("org", "", "the name of the organization you want to stop receiving notifications from")

type shortRepo struct {
	owner string
	name string
}

func main() {
	flag.Parse()
	if len(*org) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Printf("we will get rid of notifications from %s", *org)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GHACCESSTOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var unwatchQueue []shortRepo

	for pg := 1; ; pg++ {
		repos, _, err := client.Activity.ListWatched(ctx, "", &github.ListOptions{Page: pg, PerPage: 100})

		if err != nil {
			log.Fatal(err)
		}

		if len(repos) == 0 {
			break
		}
		for _, r := range repos {
			sr := shortRepo{owner: r.GetOwner().GetLogin(), name: r.GetName()}
			if sr.owner == *org {
				unwatchQueue = append(unwatchQueue, sr)
				log.Printf("queuing %+v to unwatch", sr)
			}
		}
	}
	for _, sr := range unwatchQueue {
		log.Printf("unwatching %+v...", sr)
		_, err := client.Activity.DeleteRepositorySubscription(ctx, sr.owner, sr.name)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Print("done!")
}
