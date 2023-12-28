package server

import (
	"net/http"

	"github.com/bobykurniawan11/starter-go-prisma/controllers"
	"github.com/bobykurniawan11/starter-go-prisma/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	v1 := router.Group("v1")
	{

		userGroup := v1.Group("users")
		{
			userGroup.Use(middlewares.JwtAuthMiddleware())
			user := new(controllers.UserController)
			userGroup.GET("/", user.GetAll)
			userGroup.POST("/", user.CreateUser)
			userGroup.GET("/:id", user.GetUserById)
			userGroup.PUT("/:id", user.UpdateUser)
			userGroup.DELETE("/:id", user.DeleteUser)

		}
		authGroup := v1.Group("auth")
		{
			auth := new(controllers.AuthController)
			authGroup.POST("/sign-in", auth.SignIn)
			authGroup.POST("/sign-up", auth.SignUp)
			authGroup.POST("/avatar", auth.UploadAvatar)
			authGroup.GET("/me", auth.Me)

		}

	}

	return router

}
