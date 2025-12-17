package main

import (
	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view <collection> <uuid>",
	Short: "View a single record by UUID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := getClient()
		name := args[0]
		uuid := args[1]

		sel := qdrant.WithPayloadSelector{
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true},
		}

		points, err := client.Get(ctx, &qdrant.GetPoints{
			CollectionName: name,
			Ids: []*qdrant.PointId{
				qdrant.NewID(uuid),
			},
			WithPayload: &sel,
		})
		if err != nil {
			return err
		}

		if recordsFormat == "table" {
			printPointsTable(points, recordsFields)
		} else if recordsFormat == "json" {
			output(points, recordsFormat)
		} else {
			printPointsTSV(points, recordsFields)
		}
		return nil
	},
}
