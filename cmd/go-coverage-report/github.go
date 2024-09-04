// Credit: github.com/flipgroup/golang-cover-diff
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func createOrUpdateGithubComment(ctx context.Context, summary string) {
	const commentMarker = "<!-- info:golang-coverage-report -->"

	auth_token := os.Getenv("GITHUB_TOKEN")
	if auth_token == "" {
		fmt.Println("no GITHUB_TOKEN, not reporting to GitHub.")
		return
	}

	ownerAndRepo := os.Getenv("GITHUB_REPOSITORY")
	if ownerAndRepo == "" {
		fmt.Println("no GITHUB_REPOSITORY, not reporting to GitHub.")
		return
	}

	parts := strings.SplitN(ownerAndRepo, "/", 2)
	owner := parts[0]
	repo := parts[1]

	prNumStr := os.Getenv("GITHUB_PULL_REQUEST_ID")
	if prNumStr == "" {
		fmt.Println("no GITHUB_PULL_REQUEST_ID, not reporting to GitHub.")
		return
	}

	prNum, err := strconv.Atoi(prNumStr)
	if err != nil {
		fmt.Println("provided GITHUB_PULL_REQUEST_ID is not a valid number, not reporting to GitHub.")
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: auth_token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNum, &github.IssueListCommentsOptions{})
	if err != nil {
		panic(err)
	}

	// iterate over existing pull request comments - if existing coverage comment found then update
	body := buildCommentBody(commentMarker, summary)
	for _, c := range comments {
		if c.Body == nil {
			continue
		}

		if *c.Body == body {
			// existing comment body is identical - no change
			return
		}

		if strings.HasPrefix(*c.Body, commentMarker) {
			// found existing coverage comment - update
			_, _, err := client.Issues.EditComment(ctx, owner, repo, *c.ID, &github.IssueComment{
				Body: &body,
			})
			if err != nil {
				panic(err)
			}
			return
		}
	}

	// no coverage comment found - create
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNum, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		panic(err)
	}
}

func buildCommentBody(commentMarker, summary string) string {
	return commentMarker + "\n" + summary
}
