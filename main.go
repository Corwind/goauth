package main

import (
	"github.com/Corwind/goauth/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

func UserRegistration(router *gin.RouterGroup, env *routes.Env) {
	router.POST("/registration", env.V1PostUsers)
	router.POST("/login", env.V1LoginUser)
}

func UserRetrieve(router *gin.RouterGroup, env *routes.Env) {
	router.GET("", env.V1GetUsers)
	router.GET("/:id", env.V1GetUser)
}

func OrganizationCreation(router *gin.RouterGroup, env *routes.Env) {
	router.POST("", env.V1PostOrganizations)
}

func OrganizationRetrieve(router *gin.RouterGroup, env *routes.Env) {
	router.GET("", env.V1GetOrganizations)
	router.GET("/:id", env.V1GetOrganization)
}

func SetupCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:3002"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func main() {
	fdb.MustAPIVersion(630)
	router := gin.Default()

	SetupCors(router)

	env := routes.Env{
		DB: fdb.MustOpenDefault(),
	}

	v1 := router.Group("/api/v1")
	UserRegistration(v1.Group("/"), &env)
	v1.Use(env.AuthMiddleWare(true))
	userRetrieveGroup := v1.Group("/users")
	userRetrieveGroup.Use(env.AuthorizationMiddleWare())
	UserRetrieve(userRetrieveGroup, &env)
	router.Run(":3000")
}
