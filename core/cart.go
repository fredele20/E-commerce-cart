package core

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fredele20/e-commerce-cart/database"
	"github.com/fredele20/e-commerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}


func (app *Application) AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var prodCtx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err = database.AddProductToCart(prodCtx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(200, "successfully added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var prodCtx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err = database.RemoveCartItem(prodCtx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(200, "item successfully removed from cart")
	}
}


func GetItemFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Query("id")
		if user_id == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			ctx.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var userCtx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var filledcart models.User
		err := userCollection.FindOne(userCtx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(filledcart)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(404, "not found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, 
		{Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointcursor, err := userCollection.Aggregate(userCtx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Panicln(err)
		}

		var listing []bson.M
		if err = pointcursor.All(userCtx, &listing); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			ctx.IndentedJSON(200, json["total"])
			ctx.IndentedJSON(200, filledcart.UserCart)
		}

		userCtx.Done()
	}
}


func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userQueryID := ctx.Query("id")
		if userQueryID == "" {
			log.Panic("user id is empty")
			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
			return
		}

		var userCtx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		err := database.BuyItemFromCart(userCtx, app.userCollection, userQueryID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(200, "successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userId")
		if userQueryID == "" {
			log.Println("user id is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var prodCtx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err = database.InstantBuyer(prodCtx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(200, "successfully placed the order")
	}
}