package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kzinthant-d3v/hotel-reservation/db"
	"github.com/kzinthant-d3v/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 100,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 100,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 90,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(insertedHotel)
	return nil
}

func main() {
	seedHotel("Hilton", "Yangon", 3)
	seedHotel("Shangri-La", "Yangon", 4)
	seedHotel("Hilton", "Mandalay", 5)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
