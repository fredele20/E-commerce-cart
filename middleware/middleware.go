package middleware

import (
	"net/http"

	"github.com/fredele20/e-commerce-cart/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provided"})
			ctx.Abort()
			return
		}

		claims, err := tokens.ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("uid", claims.Uid)
		ctx.Next()
	}
}
