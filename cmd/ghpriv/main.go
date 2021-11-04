package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GHACCESSTOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		panic(err)
	}
	page := 1
	var reposToTurnPrivate []*github.Repository
	for {
		repos, _, err := client.Repositories.List(ctx, user.GetLogin(), &github.RepositoryListOptions{Type: "public", ListOptions: github.ListOptions{Page: page}})
		if err != nil {
			panic(err)
		}
		if len(repos) == 0 {
			break
		}
		for _, repo := range repos {
			if repo.GetStargazersCount()+repo.GetForksCount() < 5 {
				reposToTurnPrivate = append(reposToTurnPrivate, repo)
			}
		}
		page++
	}

	for _, repo := range reposToTurnPrivate {
		owner, name := repo.GetOwner().GetLogin(), repo.GetName()
		log.Printf("https://github.com/%s/%s (fork:%v)", owner, name, repo.GetFork())
		var whatDo string
		for !(whatDo == "private" || whatDo == "delete" || whatDo == "keep") {
			fmt.Scanln(&whatDo)
			log.Printf("got:%s", whatDo)
		}
		if whatDo == "private" {
			private := true
			_, _, err := client.Repositories.Edit(ctx, owner, name, &github.Repository{Private: &private})
			if err != nil {
				panic(err)
			}
		} else if whatDo == "delete" {
			_, err = client.Repositories.Delete(ctx, owner, name)
			if err != nil {
				panic(err)
			}
		} else if whatDo == "keep" {
		}
	}

	log.Print("done!")
}
