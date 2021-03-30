package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func reseed() routing.Handler {
	return func(c *routing.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://carts-db:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("cart").Collection("data")

		var documents []interface{}
		for i := 0; i < 100000; i++ {
			json := fmt.Sprintf(`{"_class":"works.weave.socks.cart.entities.Cart","customerId":"%s","items":[]}`, uuid.New().String())
			var doc interface{}
			if err := bson.UnmarshalExtJSON([]byte(json), true, &doc); err != nil {
				return fmt.Errorf("expected bson.UnmarshalExtJSON(\njson = %s\n) returns nil err; got err = %w", json, err)
			}

			documents = append(documents, doc)
		}

		if _, err = collection.InsertMany(ctx, documents); err != nil {
			return fmt.Errorf("expected collection.InsertMany returns nil err; got err = %w", err)
		}

		fmt.Println("inserted")
		return nil
	}
}

func main() {
	router := routing.New()
	router.Post("/mode", reseed())

	if err := fasthttp.ListenAndServe(":8079", router.HandleRequest); err != nil {
		panic(fmt.Errorf("expected fasthttp.ListenAndServe() returns nil err; got err = %w", err))
	}
}
