package JobsRoute

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func JobsRoutes(adminRoute *gin.Engine) {
	//enable cors
	adminRoute.Use(cors.Default())

}
