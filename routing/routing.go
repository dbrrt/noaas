package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NoaasRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("httpOrHttps", isHttpOrHttps)
	}

	router.GET("/", new(HealthController).Status)

	v1 := router.Group("v1")

	v1.PUT("/services/:name", new(ServiceController).ServiceProvisionner)

	return router

}

func SetupServer() *gin.Engine {
	return NoaasRouter()
}
