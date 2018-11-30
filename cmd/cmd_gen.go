package cmd

import "github.com/spf13/cobra"

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "generators",
}

func init() {
	cmdRoot.AddCommand(cmdGen)
}
