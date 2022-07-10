package controllers

import (
	"userapi/config"

	"github.com/gin-gonic/gin"
)

func StartMainServer() {

	// router 設定
	r := gin.New()

	//--- handler 設定
	r.POST("/signup", signup)

	r.Run(":" + config.Config.Port)

}
