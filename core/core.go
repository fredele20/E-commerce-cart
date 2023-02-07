package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fredele20/e-commerce-cart/database"
	"github.com/fredele20/e-commerce-cart/models"
	"github.com/fredele20/e-commerce-cart/tokens"
	"github.com/fredele20/e-commerce-cart/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "Users")
var productCollection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(userCtx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}

		count, err = userCollection.CountDocuments(userCtx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with this phone number already exists"})
			return
		}

		password := utils.HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refereshToken, _ := tokens.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)

		user.Token = &token
		user.Referesh_Token = &refereshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, insertErr := userCollection.InsertOne(userCtx, user)
		if insertErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "the user was not created"})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusCreated, "successfully created a user")
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var founduser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := userCollection.FindOne(userCtx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		validPassword, msg := utils.VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()

		if !validPassword {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refereshtoken, _ := tokens.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		tokens.UpdateAllTokens(token, refereshtoken, founduser.User_ID)

		ctx.JSON(http.StatusFound, founduser)
	}

}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var productList []models.Product
		var prodCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := productCollection.Find(prodCtx, bson.D{{}})
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, "something failed.")
			return
		}

		err = cursor.All(prodCtx, &productList)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(prodCtx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "invalid")
			return
		}

		defer cancel()

		ctx.IndentedJSON(200, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var searchProduct []models.Product
		queryParam := ctx.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			ctx.Abort()
			return
		}

		var prodCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQueryDB, err := productCollection.Find(prodCtx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			ctx.IndentedJSON(404, "something failed while searching DB")
			return
		}

		err = searchQueryDB.All(prodCtx, &searchProduct)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "invalid")
			return
		}

		defer searchQueryDB.Close(prodCtx)

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(400, "invalid request")
			return
		}

		defer cancel()
		ctx.IndentedJSON(200, searchProduct)
	}
}
