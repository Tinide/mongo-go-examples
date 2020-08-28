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
 * find people who live in London and like the book "Journey to the West"
 * only displays matched.
 */
func TestAggregateArray(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M

	client = getMongoClient()
	defer client.Disconnect(ctx)
	seedFavoritesData(client, dbName)

	pipeline := `
	[{
		"$match": {
			"favoritesList": {
				"$elemMatch": {
					"city": "London",
					"book": "Journey to the West"
				}
			}
		}
	}, {
		"$project": {
			"favoritesList": {
				"$filter": {
					"input": "$favoritesList",
					"as": "favorite",
					"cond": {
						"$eq": ["$$favorite.book", "Journey to the West"]
					}
				}
			},
			"_id": 0,
			"email": 1
		}
	}, {
		"$unwind": {
			"path": "$favoritesList"
		}
	}]`
	collection = client.Database(dbName).Collection(collectionFavorites)
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, MongoPipeline(pipeline), opts); err != nil {
		t.Fatal(err)
	}
	defer cur.Close(ctx)
	total := 0
	for cur.Next(ctx) {
		cur.Decode(&doc)
		t.Log(doc["email"], "likes movie", "'", doc["favoritesList"].(bson.M)["movie"], "' too.")
		total++
	}
	t.Log("total", total)
}

func TestAggregateConcatArrays(t *testing.T) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M
	client = getMongoClient()
	defer client.Disconnect(ctx)
	seeded := seedFavoritesData(client, dbName)

	pipeline := `
	[{
		'$project': {
			'name': {
				'$concat': [
					'$firstName', ' ', '$lastName'
				]
			},
			'books': {
				'$map': {
					'input': '$favoritesKVList',
					'as': 'fa',
					'in': {
						'$filter': {
							'input': '$$fa.categories',
							'as': 'fa',
							'cond': {
								'$eq': [
									'$$fa.key', 'book'
								]
							}
						}
					}
				}
			}
		}
	}, {
		'$project': {
			'name': 1,
			'books': {
				'$reduce': {
					'input': '$books',
					'initialValue': [],
					'in': {
						'$concatArrays': [
							'$$value', '$$this'
						]
					}
				}
			}
		}
	}, {
		'$project': {
			'_id': 0,
			'name': 1,
			'books': '$books.value'
		}
	}]`

	collection = client.Database(dbName).Collection(collectionFavorites)
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, MongoPipeline(pipeline), opts); err != nil {
		t.Fatal(err)
	}
	defer cur.Close(ctx)
	total := 0
	for cur.Next(ctx) {
		cur.Decode(&doc)
		total++
	}

	if seeded != int64(total) {
		t.Fatal("expected", seeded, "but got", total)
	}
}
