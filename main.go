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
	"strconv"
	"time"
)

func clear() routing.Handler {
	return func(c *routing.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://carts-db:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("data").Collection("cart")

		if err = collection.Drop(ctx); err != nil {
			return fmt.Errorf("expected collection.Drop returns nil err; got err = %w", err)
		}

		return c.Write("cleared")
	}
}

func seed() routing.Handler {
	return func(c *routing.Context) error {
		numSeed := 100000
		if len(c.Param("num")) > 0 {
			num, err := strconv.Atoi(c.Param("num"))
			if err != nil {
				return fmt.Errorf("could not convert num param = %s to integer; err = %w", c.Param("num"), err)
			}

			if num <= 0 {
				return fmt.Errorf("num must be greater than 0; got num = %d", num)
			}

			numSeed = num
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://carts-db:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("data").Collection("cart")

		var documents []interface{}
		for i := 0; i < numSeed; i++ {
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

		return c.Write(fmt.Sprintf("inserted %d documents", numSeed))
	}
}

func main() {
	router := routing.New()
	router.Put("/db", seed())
	router.Put("/db/<num:\\d+>", seed())
	router.Delete("/db", clear())

	if err := fasthttp.ListenAndServe(":8079", router.HandleRequest); err != nil {
		panic(fmt.Errorf("expected fasthttp.ListenAndServe() returns nil err; got err = %w", err))
	}
}
