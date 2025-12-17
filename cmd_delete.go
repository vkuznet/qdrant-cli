package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <collection>",
	Short: "Delete collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		return deleteCollection(ctx, client, args[0])
	},
}

