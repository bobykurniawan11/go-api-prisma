package controllers

import (
	"net/http"

	"github.com/bobykurniawan11/starter-go-prisma/db"
	"github.com/bobykurniawan11/starter-go-prisma/utils"
	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterUser struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

type UploadAvatar struct {
	Avatar string `json:"avatar" binding:"required"`
}

type AuthController struct{}

func (u AuthController) SignIn(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	var input LoginUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}

	_existingUser, err := client.User.FindUnique(
		db.User.Email.Equals(input.Email),
	).Exec(c)

	if err != nil {
		c.JSON(422, gin.H{"error": "Email not found"})
		client.Prisma.Disconnect()
		return
	}

	if !utils.CheckPasswordHash(input.Password, _existingUser.Password) {
		c.JSON(422, gin.H{"error": "Invalid password"})
		client.Prisma.Disconnect()
		return
	}

	token, err := utils.GenerateToken(
		_existingUser.ID,
	)

	c.JSON(200, gin.H{"token": token})

}

func (u AuthController) SignUp(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	var input RegisterUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}
	_existingUser, _ := client.User.FindUnique(
		db.User.Email.Equals(input.Email),
	).Exec(c)
	if _existingUser != nil {
		c.JSON(422, gin.H{"error": "Email already exist"})
		client.Prisma.Disconnect()
		return
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	user, err := client.User.CreateOne(
		db.User.Email.Set(input.Email),
		db.User.Password.Set(hashedPassword),
		db.User.Name.Set(input.Name),
	).Exec(c)

	if err != nil {
		c.JSON(422, gin.H{"error": "Error creating user"})
		client.Prisma.Disconnect()
		return
	}
	token, err := utils.GenerateToken(
		user.ID,
	)
	c.JSON(200, gin.H{"token": token})
}

func (u AuthController) UploadAvatar(c *gin.Context) {

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	// single file
	file, _ := c.FormFile("file")

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "./uploads/"+file.Filename)

	id, err := utils.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := client.User.FindUnique(
		db.User.ID.Equals(id.String()),
	).Update(
		db.User.Avatar.Set(file.Filename),
	).Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"user": user})

}

func (u AuthController) Me(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	//make as uuid tokenString

	id, err := utils.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := client.User.FindUnique(
		db.User.ID.Equals(id.String()),
	).Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
