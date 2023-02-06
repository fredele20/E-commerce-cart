package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fredele20/e-commerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "something failed"})
			ctx.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
			return
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()
		if err = ctx.BindJSON(&addresses); err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var addrCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := userCollection.Aggregate(addrCtx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
			return
		}

		var addressInfo []bson.M
		if err = pointcursor.All(addrCtx, addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := userCollection.UpdateOne(addrCtx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			ctx.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		addrCtx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			ctx.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
			return
		}

		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var addrCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = userCollection.UpdateOne(addrCtx, filter, update)
		if err != nil {
			ctx.IndentedJSON(500, "something failed")
			return
		}

		defer cancel()
		addrCtx.Done()
		ctx.IndentedJSON(200, "successfully updated the home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
			ctx.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
			return
		}

		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}

		var addrCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.Pincode}}}}

		_, err = userCollection.UpdateOne(addrCtx, filter, update)
		if err != nil {
			ctx.IndentedJSON(500, "something failed")
			return
		}

		defer cancel()
		addrCtx.Done()
		ctx.JSON(200, "successfully updated the work address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid search index"})
			ctx.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			ctx.IndentedJSON(500, "Internal server error")
			return
		}

		var addrCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = userCollection.UpdateOne(addrCtx, filter, update)
		if err != nil {
			ctx.IndentedJSON(404, "something failed")
			return
		}

		defer cancel()
		addrCtx.Done()
		ctx.IndentedJSON(200, "successfully deleted")
	}
}
