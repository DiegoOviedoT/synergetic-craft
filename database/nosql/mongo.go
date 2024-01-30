package nosql

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientMongo interface {
	InsertOne(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter bson.D, collection string, document interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter bson.D, collection string) (*mongo.DeleteResult, error)
	Close()
}

type DocConfig struct {
	DocDBName     string
	ConnectionStr string
}

type mongoDB struct {
	docDBName string
	client    *mongo.Client
}

func NewMongo(conf DocConfig) (*mongoDB, error) {
	client, err := docDBConnect(conf)
	if err != nil {
		return nil, err
	}

	return &mongoDB{
		docDBName: conf.DocDBName,
		client:    client,
	}, nil
}

func docDBConnect(conf DocConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	options := options.Client()
	options.ApplyURI(conf.ConnectionStr)

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *mongoDB) InsertOne(ctx context.Context, docCollection string, document interface{}) (*mongo.InsertOneResult, error) {
	return m.client.Database(m.docDBName).Collection(docCollection).InsertOne(ctx, document)

}

func (m *mongoDB) UpdateOne(ctx context.Context, filter bson.D, docCollection string, document interface{}) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": document}

	return m.client.Database(m.docDBName).Collection(docCollection).UpdateOne(ctx, filter, update, opts)
}

func (m *mongoDB) DeleteOne(ctx context.Context, filter bson.D, docCollection string) (*mongo.DeleteResult, error) {
	return m.client.Database(m.docDBName).Collection(docCollection).DeleteOne(ctx, filter)
}

func (m *mongoDB) Find(ctx context.Context, docCollection string) {
	//m.client.Database(m.docDBName).Collection(docCollection).Find()
}

func (m *mongoDB) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
