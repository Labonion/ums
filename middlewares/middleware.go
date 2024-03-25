package middlewares

import (
	"markie-backend/database"
	"markie-backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey = c.GetHeader("X-API-KEY")
		token, err := database.GetData(database.RedisClient, apiKey)
		if err != nil {
			utils.Failure(c, bson.M{"message": "Could not Authenticate"}, http.StatusBadRequest)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenParts := strings.Split(token, " ")
		if tokenParts[0] != "Bearer" || tokenParts[0] == "" || tokenParts[0] == "Bearer " {
			utils.Failure(c, bson.M{"message": "Token malformed"}, http.StatusBadRequest)
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			tokenString := tokenParts[1]
			if len(tokenString) == 195 {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					userID := claims["jti"].(string)
					c.Set("userId", userID)
					c.Next()
				} else if validationError, isOk := err.(jwt.ValidationError); isOk {
					if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
						utils.Failure(c, bson.M{"message": "Token Invalid"}, http.StatusBadRequest)
						c.AbortWithStatus(http.StatusUnauthorized)

					} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						utils.Failure(c, bson.M{"message": "Token Expired"}, http.StatusBadRequest)
						c.AbortWithStatus(http.StatusUnauthorized)

					} else {
						utils.Failure(c, bson.M{"message": "Token Un-recognized"}, http.StatusBadRequest)
						c.AbortWithStatus(http.StatusUnauthorized)

					}
				} else {
					utils.Failure(c, bson.M{"message": "Token Un-recognized"}, http.StatusBadRequest)
					c.AbortWithStatus(http.StatusUnauthorized)

				}
			} else {
				utils.Failure(c, bson.M{"message": "Token Invalid"}, http.StatusBadRequest)
				c.AbortWithStatus(http.StatusUnauthorized)

			}
		}
	}
}
