package routes

import (
	"github.com/fredele20/e-commerce-cart/core"
	"github.com/gin-gonic/gin"
)


func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", core.Signup())
	incomingRoutes.POST("/users/login", core.Login())
	incomingRoutes.POST("/admin/addproduct", core.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", core.SearchProduct())
	incomingRoutes.GET("/users/search", core.SearchProductByQuery())
}