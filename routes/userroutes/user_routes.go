package userroutes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rms_api/internal/controllers/usercontrollers"
	authorization "github.com/rms_api/internal/middlewares/autherisation"
	"github.com/rms_api/routes"
)

func UserRoutes(userRoute *gin.Engine) {
	//enable cors
	userRoute.Use(cors.Default())
	userRoute.POST(routes.GenericUserLogin, func(ctx *gin.Context) {
		usercontrollers.LogInService(ctx)
	})
	userRoute.POST(routes.GenericUserSignup, func(context *gin.Context) {
		usercontrollers.SignUpService(context)
	})
	userRoute.POST(routes.ResumeUpload, authorization.AuthorizeRequest(), func(context *gin.Context) {
		usercontrollers.UploadResume(context)
	})
}
