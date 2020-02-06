package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	databaseURI := os.Getenv("ATLAS_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		cancel()
		log.Fatalf("could not connect to database server: %v:", err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("podcasts_app")
	// podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	readAll(ctx, episodesCollection)
}

func insertData(ctx context.Context, podcastsCollection, episodesCollection *mongo.Collection) {
	podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "The Polygot Developer Podcst"},
		{Key: "author", Value: "Nic Raboy"},
		{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "podcast", Value: podcastResult.InsertedID},
			{Key: "title", Value: "GraphQL for API Development"},
			{Key: "description", Value: "Learn about GraphQL from the go-creator of GraphQL, Lee Byron."},
			{Key: "duration", Value: 25},
		},
		bson.D{
			{Key: "podcast", Value: podcastResult.InsertedID},
			{Key: "title", Value: "Progressive Web Application Development"},
			{Key: "description", Value: "Learn about PWA development with Tara Manicsic."},
			{Key: "duration", Value: 32},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Inserted %v documents into episode collection!\n", len(episodeResult.InsertedIDs))
}

func readAll(ctx context.Context, collection *mongo.Collection) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// var results []bson.M
	// if err = cursor.All(ctx, &results); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(results)

	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}
}
