package commands

import (
	"errors"
	"fmt"

	"github.com/MichaelMure/git-bug/bug"
	"github.com/MichaelMure/git-bug/cache"
	"github.com/spf13/cobra"
)

func runPull(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Only pulling from one remote at a time is supported")
	}

	remote := "origin"
	if len(args) == 1 {
		remote = args[0]
	}

	backend, err := cache.NewRepoCache(repo)
	if err != nil {
		return err
	}
	defer backend.Close()

	fmt.Println("Fetching remote ...")

	stdout, err := backend.Fetch(remote)
	if err != nil {
		return err
	}

	fmt.Println(stdout)

	fmt.Println("Merging data ...")

	for merge := range backend.MergeAll(remote) {
		if merge.Err != nil {
			return merge.Err
		}

		if merge.Status != bug.MergeStatusNothing {
			fmt.Printf("%s: %s\n", merge.Bug.HumanId(), merge.Status)
		}
	}

	return nil
}

// showCmd defines the "push" subcommand.
var pullCmd = &cobra.Command{
	Use:   "pull [<remote>]",
	Short: "Pull bugs update from a git remote",
	RunE:  runPull,
}

func init() {
	RootCmd.AddCommand(pullCmd)
}
