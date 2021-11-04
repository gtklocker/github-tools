package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

var org = flag.String("org", "", "the name of the organization you want to filter repos from")
var unwatch = flag.Bool("unwatch", false, "whether to unwatch all filtered repos")
var unstar = flag.Bool("unstar", false, "whether to unstar all filtered repos")

type shortRepo struct {
	owner string
	name  string
}

func toShortRepo(r *github.Repository) *shortRepo {
	return &shortRepo{owner: r.GetOwner().GetLogin(), name: r.GetName()}
}

func doUnwatch(ctx context.Context, cl *github.Client) {
	var q []*shortRepo

	for pg := 1; ; pg++ {
		repos, _, err := cl.Activity.ListWatched(ctx, "", &github.ListOptions{Page: pg, PerPage: 100})

		if err != nil {
			log.Fatal(err)
		}

		if len(repos) == 0 {
			break
		}
		for _, r := range repos {
			sr := toShortRepo(r)
			if sr.owner == *org {
				q = append(q, sr)
				log.Printf("queuing %+v to unwatch", sr)
			}
		}
	}
	for _, sr := range q {
		log.Printf("unwatching %+v...", sr)
		_, err := cl.Activity.DeleteRepositorySubscription(ctx, sr.owner, sr.name)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doUnstar(ctx context.Context, cl *github.Client) {
	var q []*shortRepo

	for pg := 1; ; pg++ {
		starredRepos, _, err := cl.Activity.ListStarred(ctx, "", &github.ActivityListStarredOptions{ListOptions: github.ListOptions{Page: pg, PerPage: 100}})

		if err != nil {
			log.Fatal(err)
		}

		if len(starredRepos) == 0 {
			break
		}
		for _, r := range starredRepos {
			sr := toShortRepo(r.GetRepository())
			if sr.owner == *org {
				q = append(q, sr)
				log.Printf("queuing %+v to unstar", sr)
			}
		}
	}
	for _, sr := range q {
		log.Printf("unstarring %+v...", sr)
		_, err := cl.Activity.Unstar(ctx, sr.owner, sr.name)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	flag.Parse()
	if *org == "" || !(*unwatch || *unstar) {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Printf("org:%s, unwatch:%v, unstar:%v", *org, *unwatch, *unstar)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GHACCESSTOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	if *unwatch {
		doUnwatch(ctx, client)
	}
	if *unstar {
		doUnstar(ctx, client)
	}
	log.Print("done!")
}
