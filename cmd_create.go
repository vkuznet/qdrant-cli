package main

import "github.com/spf13/cobra"

var createDim uint64

func init() {
	createCmd.Flags().Uint64Var(
		&createDim,
		"dim",
		50176,
		"Vector dimension",
	)

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create <collection>",
	Short: "Create collection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		return createCollection(ctx, client, args[0], createDim)
	},
}

