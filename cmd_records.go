package main

import "github.com/spf13/cobra"

var (
	recordsLimit uint
)

func init() {
	recordsCmd.Flags().UintVar(
		&recordsLimit,
		"limit",
		5,
		"Scroll page size",
	)

	rootCmd.AddCommand(recordsCmd)
}

var recordsCmd = &cobra.Command{
	Use:   "records <collection>",
	Short: "List records in a collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		name := args[0]
		return scrollCollection(ctx, client, name, recordsFields, recordsFormat, recordsFilter, recordsLimit)
	},
}
