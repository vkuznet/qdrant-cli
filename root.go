package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "qdrant-cli",
		Short: "Qdrant inspection and management CLI",
	}
	qdrantURL     string
	recordsFormat string
	recordsFields string
	recordsFilter string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// define sub-options applied to all commands
	rootCmd.PersistentFlags().StringVar(
		&qdrantURL,
		"url",
		"http://localhost:6334",
		"Qdrant gRPC URL (e.g. http://localhost:6334)",
	)

	rootCmd.PersistentFlags().StringVar(
		&recordsFormat,
		"format",
		"tsv",
		"Output format: tsv, table, json",
	)
	rootCmd.PersistentFlags().StringVar(
		&recordsFields,
		"fields",
		"id,filename",
		"Comma-separated payload fields",
	)

	rootCmd.PersistentFlags().StringVar(
		&recordsFilter,
		"filter",
		"",
		"Payload filter key=value",
	)

}

func getClient() (*qdrant.Client, context.Context) {
	u, err := url.Parse(qdrantURL)
	if err != nil {
		log.Fatal(err)
	}

	host := u.Hostname()
	portStr := u.Port()
	if portStr == "" {
		portStr = "6334"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}

	return client, context.Background()
}
