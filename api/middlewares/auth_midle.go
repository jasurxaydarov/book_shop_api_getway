package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jasurxaydarov/book_shop_api_getway/token"
)

func AuthMiddlewareUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("authorization")

		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "authorization token not provided"})
			ctx.Abort()
		}

		claim, err := token.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			ctx.Abort()
		}

		if claim.UserRole != "user" {
			ctx.JSON(401, gin.H{"error": "your role isn't user "})
			ctx.Abort()
		}

		ctx.Next()
	}
}



func AuthMiddlewareAdmin() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("authorization")

		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "authorization token not provided"})
			ctx.Abort()
		}

		claim, err := token.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			ctx.Abort()
		}

		if claim.UserRole != "admin" {
			ctx.JSON(401, gin.H{"error": "your role isn't admin "})
			ctx.Abort()
		}

		ctx.Next()
	}
}