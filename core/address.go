package core

import (
	"context"
	"net/http"
	"time"

	"github.com/fredele20/e-commerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
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
