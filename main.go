package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/qdrant/go-client/qdrant"
)

var (
	host = flag.String("host", "localhost", "Qdrant host")
	port = flag.Int("port", 6334, "Qdrant gRPC port")

	listFlag   = flag.Bool("list", false, "List collections")
	createName = flag.String("create", "", "Create collection")
	deleteName = flag.String("delete", "", "Delete collection")
	scrollName = flag.String("scroll", "", "Scroll records")
	infoName   = flag.String("info", "", "Describe collection")

	dim   = flag.Uint64("dim", 50176, "Vector dimension")
	limit = flag.Uint("limit", 5, "Scroll page size")

	filterKV = flag.String("filter", "", "Filter payload.key=value")
	format   = flag.String("format", "tsv", "format to use: tsv, json, table")

	fields = flag.String(
		"fields",
		"id,filename",
		"Comma-separated payload fields to display (default: id,filename)",
	)
)

func getClient() (*qdrant.Client, error) {
	return qdrant.NewClient(&qdrant.Config{
		Host: *host,
		Port: *port,
	})
}

func main() {
	flag.Parse()

	client, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	switch {
	case *listFlag:
		err = listCollections(ctx, client)

	case *infoName != "":
		err = describeCollection(ctx, client, *infoName)

	case *createName != "":
		err = createCollection(ctx, client, *createName, *dim)
		/*
			if err == nil {
				err = insertDummyData(ctx, client, *createName, *dim)
			}
		*/

	case *deleteName != "":
		err = deleteCollection(ctx, client, *deleteName)

	case *scrollName != "":
		err = scrollCollection(ctx, client, *scrollName, *limit)

	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}
