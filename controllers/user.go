package controllers

import (
	"fmt"
	"net/http"

	"github.com/bobykurniawan11/starter-go-prisma/db"
	"github.com/bobykurniawan11/starter-go-prisma/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct{}
type InputUser struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUser struct {
	//NAme Optional
	Name string `json:"name"`
	//Email Optional

	Phone string `json:"phone"`
}

func (u UserController) GetAll(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	users, err := client.User.FindMany().Exec(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (u UserController) CreateUser(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	var input InputUser
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	fmt.Println(input.Email, input.Name, input.Password, hashedPassword)

	user, err := client.User.CreateOne(
		db.User.Email.Set(input.Email),
		db.User.Password.Set(hashedPassword),
		db.User.Name.Set(input.Name),
	).Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (u UserController) GetUserById(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	id := c.Param("id")
	user, err := client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
func (u UserController) UpdateUser(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	id := c.Param("id")

	// Define a struct to hold the fields you want to update
	var updateFields UpdateUser

	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the existing user from the database
	_user, err := client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.Name.Set(updateFields.Name),
		db.User.Phone.Set(updateFields.Phone),
	).Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": _user})

}

func (u UserController) DeleteUser(c *gin.Context) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	id := c.Param("id")

	User, err := client.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete user", "user": User})
}
