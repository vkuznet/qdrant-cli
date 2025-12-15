package main

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func createCollection(ctx context.Context, c *qdrant.Client, name string, dim uint64) error {
	fmt.Printf("creating %s collection with size %d\n", name, dim)
	params := &qdrant.VectorParams{
		Size:     uint64(dim),
		Distance: qdrant.Distance_Cosine,
	}
	colClient := c.GetCollectionsClient()
	_, err := colClient.Create(ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig:  qdrant.NewVectorsConfig(params),
	})
	return err
}

func deleteCollection(ctx context.Context, c *qdrant.Client, name string) error {
	return c.DeleteCollection(ctx, name)
}

/*
func insertDummyData(ctx context.Context, c *qdrant.Client, name string, dim uint64) error {

	vec := make([]float32, dim)
	for j := range vec {
		vec[j] = rand.Float32()
	}

	pstr := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(uint64(1)),
		Vectors: qdrant.NewVectors(vec...),
		Payload: map[string]*qdrant.Value{
			"type": qdrant.NewValueString("dummy"),
			"idx":  qdrant.NewValueInt(int64(1)),
		},
	}
	points := &qdrant.UpsertPoints{
		CollectionName: name,
		Points:         []*qdrant.PointStruct{pstr},
	}

	_, err := c.Upsert(ctx, points)
	return err
}
*/
