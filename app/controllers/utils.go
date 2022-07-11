package controllers

import (
	"userapi/config"
)

var deployEnv = config.Config.Deploy
var serverPort = config.Config.Port
