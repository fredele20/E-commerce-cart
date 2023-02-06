package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	searchFromDB, err := productCollection.Find(ctx, bson.M{"_id": productId})
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

func RemoveCartItem(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = productCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveCartItem
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// fetch the cart of the user
	// find the cart total
	// create an order with the items
	// added order to the user collection
	// added items in the cart to order list
	// empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Order_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key:"path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err = currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var total_price int32

	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	orderCart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, productCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	var productDetails models.ProductUser
	var ordersDetail models.Order

	ordersDetail.Order_ID = primitive.NewObjectID()
	ordersDetail.Order_At = time.Now()
	ordersDetail.Order_Cart = make([]models.ProductUser, 0)
	ordersDetail.Payment_Method.COD = true

	err = productCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}

	ordersDetail.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordersDetail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 :=bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	return nil
}
