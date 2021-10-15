package routes

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Corwind/goauth/users"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func stripBearerPrefixFromTokenString(token string) (string, error) {
	if len(token) > 5 && strings.ToUpper(token[0:6]) == "TOKEN " {
		return token[6:], nil
	}
	return token, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	stripBearerPrefixFromTokenString,
}

var AuthExtractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func GenerateToken(id string) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token, _ := jwt_token.SignedString([]byte("super secret key"))
	return token
}

func UpdateContextWithUser(c *gin.Context, db fdb.Database, user_id string) {
	user_, err := users.FetchUserById(db, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
		return
	}

	c.Set("user_id", user_id)
	c.Set("user", user_)
}

func (env *Env) AuthMiddleWare(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, AuthExtractor, func(token *jwt.Token) (interface{}, error) {
			return ([]byte("super secret key")), nil
		})

		if auto401 && (err != nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user_id := claims["id"].(string)
			UpdateContextWithUser(c, env.DB, user_id)
		}
	}
}

func (env *Env) AuthorizationMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id != "" {
			user := c.MustGet("user")
			user_ := user.(users.User)
			if (&user_).Id != id {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Unauthorized"))
				return
			}
		}
		c.Next()
	}
}
