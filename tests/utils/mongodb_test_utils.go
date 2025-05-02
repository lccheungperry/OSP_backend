package utils

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DefaultMongoURI = "mongodb://localhost:27017"

func GetMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	return DefaultMongoURI
}

func SetupTestDB(t *testing.T, dbName string) (*mongo.Database, func()) {
	uri := GetMongoURI()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	require.NoError(t, err)

	db := client.Database(dbName)

	cleanup := func() {
		err := db.Drop(context.Background())
		require.NoError(t, err)
		err = client.Disconnect(context.Background())
		require.NoError(t, err)
	}

	return db, cleanup
}

func CleanCollection(t *testing.T, db *mongo.Database, collectionName string) {
	err := db.Collection(collectionName).Drop(context.Background())
	require.NoError(t, err)
}
