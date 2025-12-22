package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "records <collection>",
	Short: "List records in a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		name := args[0]
		return infoCollection(ctx, client, name)
	},
}
