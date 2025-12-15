package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/qdrant/go-client/qdrant"
)

func output(v any) {
	if *format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(v)
		return
	}
	fmt.Printf("%+v\n", v)
}

func buildFilter() *qdrant.Filter {
	if *filterKV == "" {
		return nil
	}

	parts := strings.SplitN(*filterKV, "=", 2)
	if len(parts) != 2 {
		log.Fatal("filter must be payload.key=value")
	}

	return &qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewMatchKeyword(parts[0], parts[1]),
		},
	}
}

func printCollectionsTable(cols []string) {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(w, "COLLECTION")
	fmt.Fprintln(w, "----------")
	for _, c := range cols {
		fmt.Fprintf(w, "%s\n", c)
	}
	w.Flush()
}

func listCollections(ctx context.Context, c *qdrant.Client) error {
	cols, err := c.ListCollections(ctx)
	if err != nil {
		return err
	}

	if *format == "json" {
		output(cols)
	} else {
		printCollectionsTable(cols)
	}
	return nil
}

func printCollectionInfo(info *qdrant.CollectionInfo) {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)

	fmt.Fprintln(w, "FIELD\tVALUE")
	fmt.Fprintln(w, "-----\t-----")

	fmt.Fprintf(w, "Status\t%v\n", info.Status)
	fmt.Fprintf(w, "Segments\t%d\n", info.SegmentsCount)
	fmt.Fprintf(w, "Points\t%d\n", *info.PointsCount)
	fmt.Fprintf(w, "Indexed vectors\t%d\n", *info.IndexedVectorsCount)

	if p := info.Config.GetParams(); p != nil {
		fmt.Fprintf(w, "Shard number\t%d\n", p.ShardNumber)
		fmt.Fprintf(w, "Replication factor\t%d\n", *p.ReplicationFactor)
		fmt.Fprintf(w, "Write consistency\t%d\n", *p.WriteConsistencyFactor)
		fmt.Fprintf(w, "On-disk payload\t%v\n", p.OnDiskPayload)

		if vc := p.VectorsConfig.GetParams(); vc != nil {
			fmt.Fprintf(w, "Vector size\t%d\n", vc.Size)
			fmt.Fprintf(w, "Distance\t%s\n", vc.Distance.String())
		}
	}

	w.Flush()
}

func describeCollection(ctx context.Context, c *qdrant.Client, name string) error {
	info, err := c.GetCollectionInfo(ctx, name)
	if err != nil {
		return err
	}

	if *format == "json" {
		output(info)
	} else {
		printCollectionInfo(info)
	}
	return nil
}

func formatPointID(id *qdrant.PointId) string {
	if id == nil {
		return "-"
	}
	if id.GetUuid() != "" {
		return id.GetUuid()
	}
	return fmt.Sprintf("%d", id.GetNum())
}

func payloadValue(v *qdrant.Value) string {
	switch x := v.Kind.(type) {
	case *qdrant.Value_StringValue:
		return x.StringValue
	case *qdrant.Value_IntegerValue:
		return fmt.Sprintf("%d", x.IntegerValue)
	case *qdrant.Value_DoubleValue:
		return fmt.Sprintf("%f", x.DoubleValue)
	case *qdrant.Value_BoolValue:
		return fmt.Sprintf("%v", x.BoolValue)
	default:
		return "-"
	}
}

func printPointsTable(points []*qdrant.RetrievedPoint) {
	reqFields := parseFields()

	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tFIELD\tVALUE")
	fmt.Fprintln(w, "--\t-----\t-----")

	for _, p := range points {
		id := formatPointID(p.Id)

		first := true
		for _, field := range reqFields {
			var value string

			switch field {
			case "id":
				value = id
			default:
				v, ok := p.Payload[field]
				if !ok {
					continue
				}
				value = payloadValue(v)
			}

			if first {
				fmt.Fprintf(w, "%s\t%s\t%s\n", id, field, value)
				first = false
			} else {
				fmt.Fprintf(w, "\t%s\t%s\n", field, value)
			}
		}
	}

	w.Flush()
}

func scrollCollection(ctx context.Context, c *qdrant.Client, name string, limit uint) error {
	var offset *qdrant.PointId

	sel := qdrant.WithPayloadSelector{
		SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true},
	}
	ulim := uint32(limit)

	points, err := c.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: name,
		Limit:          &ulim,
		Offset:         offset,
		Filter:         buildFilter(),
		WithPayload:    &sel,
	})
	if err != nil {
		return err
	}

	if *format == "table" {
		printPointsTable(points)
	} else if *format == "json" {
		output(points)
	} else {
		printPointsTSV(points)
	}
	return nil
}

func parseFields() []string {
	parts := strings.Split(*fields, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func printPointsTSV(points []*qdrant.RetrievedPoint) {
	fields := parseFields()

	for _, p := range points {
		var cols []string

		for _, f := range fields {
			switch f {
			case "id":
				id := formatPointID(p.Id)
				// OPTIONAL: skip numeric IDs
				if p.Id.GetNum() != 0 {
					continue
				}
				cols = append(cols, id)

			default:
				v, ok := p.Payload[f]
				if !ok {
					cols = append(cols, "")
					continue
				}
				cols = append(cols, payloadValue(v))
			}
		}

		if len(cols) > 0 {
			fmt.Println(strings.Join(cols, "\t"))
		}
	}
}
