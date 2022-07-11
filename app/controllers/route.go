package controllers

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"userapi/app/models"

	"github.com/gin-gonic/gin"
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

type encryptPassword struct {
	PassWord string `json:"password"`
}

func createUser(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ユーザ登録")
	defer span.End()

	var json signupRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := models.GetUserByEmail(c, json.Email)
	if user.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"error_code": "その Email はすでに存在しております",
		})
	} else {
		user := models.User{
			Name:     json.Name,
			Email:    json.Email,
			PassWord: json.PassWord,
		}
		if err := user.CreateUser(c); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"Name":  json.Name,
			"Email": json.Email,
		})
	}
}

func getUserByEmail(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ユーザ参照")
	defer span.End()

	var json signupRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := models.GetUserByEmail(c, json.Email)

	c.JSON(http.StatusOK, gin.H{
		"ID":        user.ID,
		"UUID":      user.UUID,
		"Name":      user.Name,
		"Email":     user.Email,
		"PassWord":  user.PassWord,
		"CreatedAt": user.CreatedAt,
	})

}

func Encrypt(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "PW 暗号化")
	defer span.End()

	var json encryptPassword
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plaintext := json.PassWord
	cryptext := fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))

	c.JSON(http.StatusOK, gin.H{
		"PassWord": cryptext,
	})
}
