package main

import (
	"context"
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
	page := 1
	var following []*github.User
	for {
		users, _, err := client.Users.ListFollowing(ctx, "", &github.ListOptions{Page: page})
		if err != nil {
			panic(err)
		}
		if len(users) == 0 {
			break
		}
		for _, user := range users {
			following = append(following, user)
		}
		page++
	}

	for _, user := range following {
		log.Printf("unfollowing %s", user.GetLogin())
		_, err := client.Users.Unfollow(ctx, user.GetLogin())
		if err != nil {
			panic(err)
		}
	}

	log.Print("done!")
}
