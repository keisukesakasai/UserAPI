package controllers

import (
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

func signup(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ユーザ登録")
	defer span.End()

	var json signupRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exist, _ := models.GetUserByEmail(c, json.Email)
	if exist.ID != 0 {
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
