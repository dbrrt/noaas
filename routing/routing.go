package routing

import (
	"github.com/gin-gonic/gin"
)

func NoaasRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", new(HealthController).Status)

	v1 := router.Group("v1")

	v1.GET("/", new(HealthController).Status)

	return router

}

func SetupServer() *gin.Engine {
	return NoaasRouter()
}
