// Copyright 2018 Kuei-chun Chen. All rights reserved.

package examples

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
 * count vehicles by style and display dealer name, all brands and a total count of each style
 */
func TestAggregateLookup(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M

	client = getMongoClient()
	defer client.Disconnect(ctx)
	total := seedCarsData(client, dbName)

	pipeline := `
	[{
		"$group": {
			"_id": {
				"dealer": "$dealer",
				"style": "$style"
			},
			"brand": {
				"$addToSet": "$brand"
			},
			"count": {
				"$sum": 1
			}
		}
	}, {
		"$lookup": {
			"as": "dealer",
			"from": "dealers",
			"let": {
				"dealerId": "$_id.dealer"
			},
			"pipeline": [{
				"$match": {
					"$expr": {
						"$eq": [
							"$_id", "$$dealerId"
						]
					}
				}
			}, {
				"$project": {
					"_id": 0,
					"name": 1
				}
			}]
		}
	}]`
	collection = client.Database(dbName).Collection(collectionName)
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, MongoPipeline(pipeline), opts); err != nil {
		t.Fatal(err)
	}
	defer cur.Close(ctx)
	count := int64(0)
	for cur.Next(ctx) {
		cur.Decode(&doc)
		count += toInt64(doc["count"])
	}

	if total != count {
		t.Fatal("expected", total, "but got", count)
	}
}
