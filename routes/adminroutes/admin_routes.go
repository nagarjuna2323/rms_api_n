package adminroutes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(adminRoute *gin.Engine) {
	//enable cors
	adminRoute.Use(cors.Default())

}
