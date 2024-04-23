package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rms_api/internal/config/database"
	L "github.com/rms_api/internal/middlewares/logger"
	"github.com/rms_api/routes/userroutes"
)

func main() {
	port := ":3001"
	database.OpenDbConnection()
	userRouter := gin.Default()
	userroutes.UserRoutes(userRouter)
	err := userRouter.Run(port)
	if err != nil {
		L.RMSLog("E", "Unable to Serve Api on "+L.PrintStruct(port), nil)
		return
	}

}
