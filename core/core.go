package core

import "github.com/gin-gonic/gin"


func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func Login() gin.HandlerFunc {}


func ProductViewerAdmin() gin.HandlerFunc {}

func SearchProduct() gin.HandlerFunc {}

func SearchProductByQuery() gin.HandlerFunc {}