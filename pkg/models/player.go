package models

import "time"

type Item struct {
	ID       string `bson:"id"`
	Quantity int    `bson:"quantity"`
}

type Player struct {
	ID          string    `bson:"_id"`
	Username    string    `bson:"username"`
	X           float64   `bson:"x"`
	Y           float64   `bson:"y"`
	Attack      int       `bson:"attack"`
	Ranged      int       `bson:"ranged"`
	Hitpoints   int       `bson:"hitpoints"`
	Inventory   []Item    `bson:"inventory"`
	LastUpdated time.Time `bson:"last_updated"`
}
