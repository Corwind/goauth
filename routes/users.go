package routes

import (
	"net/http"

	"github.com/Corwind/goauth/users"
	validator "github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserPostSchema struct {
	Email           string `form:"email" json:"email" binding:"required" validate:"email"`
	Username        string `form:"username" json:"username" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required" validate:"min=8"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required" validate:"eqField=Password"`
}

type LoginPostSchema struct {
	Username string `form:"username" json:"username" binding:"required" validate:"username"`
	Password string `form:"password" json:"password" binding:"required" validate:"min=8"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type UserSerializer struct {
	c *gin.Context
}

func (serializer *UserSerializer) Response() LoginResponse {
	user := serializer.c.MustGet("user").(users.User)
	return LoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    GenerateToken(user.Id),
	}
}

func UserCreationValidator(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userPostSchema UserPostSchema
		if err := c.Bind(&userPostSchema); err == nil {
			validate := validator.New()
			if err := validate.Struct(&userPostSchema); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func (env *Env) V1GetUsers(c *gin.Context) {
	ret, err := users.FetchUsers(env.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &ret)
}

func (env *Env) V1PostUsers(c *gin.Context) {
	var form UserPostSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	user := *users.NewUser(
		form.Email,
		form.Username,
		string(hash),
	)
	_, err = users.SaveUser(env.DB, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &user)
}

func (env *Env) V1GetUser(c *gin.Context) {
	user_id := c.Param("id")
	user_, err := users.FetchUserById(env.DB, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, &user_)
}

func (env *Env) V1LoginUser(c *gin.Context) {
	var loginPostSchema LoginPostSchema

	c.Bind(&loginPostSchema)
	user, err := users.FetchUserByUsername(env.DB, loginPostSchema.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	user_ := user.(users.User)

	if (&user_).CheckPassword(loginPostSchema.Password) != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}
	UpdateContextWithUser(c, env.DB, user_.Id)
	var userSerializer UserSerializer = UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": userSerializer.Response()})
}
