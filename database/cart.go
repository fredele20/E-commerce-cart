package database

import (
	"context"
	"errors"
	"log"

	"github.com/fredele20/e-commerce-cart/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrProductDecodingFailed = errors.New("product decoding failed")
	ErrUserIdNotValid        = errors.New("this user is not valid")
	ErrCantRemoveCartItem    = errors.New("cannot remove this item from the cart")
	ErrCantGetItem           = errors.New("unable to get cart item")
	ErrCantBuyCartItem       = errors.New("cannot update the purchase")
	ErrCantUpdateUser        = errors.New("cannot add this product to the cart")
)

func AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	searchFromDB, err := productCollection.Find(ctx, bson.M{"_id":productId})
	if err != nil {
		log.Println(err)
		return ErrProductNotFound
	}

	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrProductDecodingFailed
	}

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem() {}

func BuyItemFromCart() {}

func InstantBuyer() {}
