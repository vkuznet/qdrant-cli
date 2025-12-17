package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List collections",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		return listCollections(ctx, client, recordsFormat)
	},
}
