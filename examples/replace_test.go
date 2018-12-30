// Copyright 2018 Kuei-chun Chen. All rights reserved.

package examples

import (
	"context"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func TestReplaceOne(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var ctx = context.Background()
	var doc = bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta"}
	var result *mongo.UpdateResult
	client = getMongoClient()
	collection = client.Database(dbName).Collection(collectionExamples)
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		t.Fatal(err)
	}
	doc["year"] = 1998
	if result, err = collection.ReplaceOne(ctx, bson.M{"_id": doc["_id"]}, doc); err != nil {
		t.Fatal(err)
	}
	if result.MatchedCount != 1 || result.ModifiedCount != 1 {
		t.Fatal("replace failed, expected 1 but got", result.MatchedCount)
	}
	res, _ := collection.DeleteMany(ctx, bson.M{"hometown": "Atlanta"})
	if res.DeletedCount != 1 {
		t.Fatal("replace failed, expected 1 but got", res.DeletedCount)
	}
}

func TestReplaceLoop(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur mongo.Cursor
	var result *mongo.UpdateResult
	var ctx = context.Background()
	var docs []interface{}
	docs = append(docs, bson.M{"hometown": "Atlanta", "year": 1998})
	docs = append(docs, bson.M{"hometown": "Jacksonville", "year": 1990})
	client = getMongoClient()
	collection = client.Database(dbName).Collection(collectionExamples)
	if _, err = collection.InsertMany(ctx, docs); err != nil {
		t.Fatal(err)
	}
	if cur, err = collection.Find(ctx, bson.M{"hometown": bson.M{"$exists": 1}}); err != nil {
		t.Fatal(err)
	}
	var doc bson.M
	for cur.Next(ctx) {
		cur.Decode(&doc)
		doc["updated"] = time.Now()
		if result, err = collection.ReplaceOne(ctx, bson.M{"_id": doc["_id"]}, doc); err != nil {
			t.Fatal(err)
		}
		if result.MatchedCount != 1 || result.ModifiedCount != 1 {
			t.Fatal("replace failed, expected 1 but got", result.MatchedCount)
		}
	}
	res, _ := collection.DeleteMany(ctx, bson.M{"hometown": bson.M{"$exists": 1}})
	if res.DeletedCount != int64(len(docs)) {
		t.Fatal("replace failed, expected", len(docs), "but got", res.DeletedCount)
	}
}