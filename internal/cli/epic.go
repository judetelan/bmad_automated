package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"bmad-automate/internal/lifecycle"
	"bmad-automate/internal/router"
)

func newEpicCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "epic <epic-id>",
		Short: "Run full lifecycle for all stories in an epic",
		Long: `Run the complete lifecycle for all stories in an epic to completion.

Finds all stories matching the pattern {epic-id}-{N}-* where N is numeric,
sorts them by story number, and runs each to completion before moving to the next.

For each story, executes all remaining workflows based on its current status:
  - backlog       → create-story → dev-story → code-review → git-commit → done
  - ready-for-dev → dev-story → code-review → git-commit → done
  - in-progress   → dev-story → code-review → git-commit → done
  - review        → code-review → git-commit → done
  - done          → skipped (story already complete)

The epic command stops on the first failure. Done stories are skipped and do not cause failure.
Status is updated in sprint-status.yaml after each successful workflow.

Example:
  bmad-automate epic 6
  # Runs 6-1-*, 6-2-*, 6-3-*, etc. each to completion in order`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			epicID := args[0]

			// Get all stories for this epic
			storyKeys, err := app.StatusReader.GetEpicStories(epicID)
			if err != nil {
				cmd.SilenceUsage = true
				return NewExitError(1)
			}

			// Create lifecycle executor with app dependencies
			executor := lifecycle.NewExecutor(app.Runner, app.StatusReader, app.StatusWriter)

			// Execute full lifecycle for each story in order
			for _, storyKey := range storyKeys {
				err := executor.Execute(ctx, storyKey)
				if err != nil {
					cmd.SilenceUsage = true
					if errors.Is(err, router.ErrStoryComplete) {
						fmt.Printf("Story %s is already complete, skipping\n", storyKey)
						continue
					}
					fmt.Printf("Error running lifecycle for story %s: %v\n", storyKey, err)
					return NewExitError(1)
				}
				fmt.Printf("Story %s completed successfully\n", storyKey)
			}

			fmt.Printf("All %d stories processed\n", len(storyKeys))
			return nil
		},
	}
}
