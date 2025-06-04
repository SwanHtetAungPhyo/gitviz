package main

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	// Open the current directory as a Git repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening git repo: %v", err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		log.Fatalf("Error getting HEAD: %v", err)
	}

	// Get the commit history starting from HEAD
	iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatalf("Error reading log: %v", err)
	}

	fmt.Println("Commit history:")
	// Iterate through commits and print hash and message
	err = iter.ForEach(func(c *object.Commit) error {
		fmt.Printf("- %s %s\n", c.Hash.String()[:7], c.Message)
		return nil
	})

	if err != nil {
		log.Fatalf("Error iterating commits: %v", err)
	}
}
