package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceController struct{}

type NewServiceRequest struct {
	Uri int `json:"uri"  binding:"required,url,http|https"`
}

func (h ServiceController) ServiceProvisionner(c *gin.Context) {

	var payload NewServiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"url":   nil,
			"error": err.Error(), // will return EOF if body is empty
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"url":   "https://www.google.com",
			"error": nil,
		})

	}

}
