//go:build integration
// +build integration

package nosql_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"synergetic-craft/database/nosql"
	"testing"
)

func TestMongoDB_InsertOne(t *testing.T) {
	t.Run("should return success when mongo is connect and insert success", func(t *testing.T) {
		f := setupMongoFixture()
		defer f.resource.Close()

		document := struct {
			Name   string `bson:"name" json:"name"`
			Status string `bson:"status" json:"status"`
		}{
			Name:   "document test",
			Status: "active",
		}

		collection := "prueba"

		result, err := f.resource.InsertOne(context.Background(), collection, document)

		assert.NoError(t, err)
		assert.NotNil(t, result.InsertedID)
	})

}

type mongoFixture struct {
	resource nosql.ClientMongo
}

func setupMongoFixture() *mongoFixture {
	config := nosql.DocConfig{
		DocDBName:     "testDB",
		ConnectionStr: "mongodb://test:passwordtest@localhost:27017/testDB",
	}

	newMongo, err := nosql.NewMongo(config)
	if err != nil {
		panic(err)
	}

	return &mongoFixture{
		resource: newMongo,
	}
}
