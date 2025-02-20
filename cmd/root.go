package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vitorfhc/webdiffer/pkg/store/jsonstore"
	"github.com/vitorfhc/webdiffer/pkg/types"
	"github.com/vitorfhc/webdiffer/pkg/webwatcher"
)

const (
	storeFile = "store.json"
)

var rootCmd = &cobra.Command{
	Use:   "webdiff",
	Short: "Webdiff is a tool to monitor changes in web pages",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new target to monitor",
		RunE:  addCmdRun,
	}

	addCmd.Flags().StringP("url", "u", "", "URL to monitor")
	addCmd.MarkFlagRequired("url")

	rootCmd.AddCommand(addCmd)

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the diffing tool",
		RunE:  runCmdRun,
	}

	rootCmd.AddCommand(runCmd)
}

func addCmdRun(cmd *cobra.Command, args []string) error {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}

	target := types.Target{
		URL: url,
	}

	store := jsonstore.NewJSONStore(storeFile)
	return store.InsertTarget(target)
}

func runCmdRun(cmd *cobra.Command, args []string) error {
	store := jsonstore.NewJSONStore(storeFile)
	w := webwatcher.NewWebWatcher(store)
	diffs, err := w.Run()
	if err != nil {
		return err
	}

	for _, diff := range diffs {
		fmt.Printf("Diff found in %q\n", diff.Target.URL)
	}

	return nil
}
