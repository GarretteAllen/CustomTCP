package database

import (
	"context"
	"customtcp/pkg/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var playerCollection *mongo.Collection

func Connect(uri string, dbName string) {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB:", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed:", err)
	}

	playerCollection = client.Database(dbName).Collection("players")
	fmt.Println("Connected to MongoDB")
}

func SavePlayerData(player models.Player) error {
	filter := bson.M{"username": player.Username}
	update := bson.M{
		"$set": player,
	}
	opts := options.Update().SetUpsert(true)
	_, err := playerCollection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetPlayerData(username string) (models.Player, error) {
	var player models.Player
	filter := bson.M{"username": username}
	err := playerCollection.FindOne(context.Background(), filter).Decode(&player)
	if err != nil {
		return models.Player{}, err
	}
	return player, nil
}
