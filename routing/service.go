package routing

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ServiceController struct{}

type NewServiceRequest struct {
	Url    string `json:"url" binding:"required,url"`
	Script bool   `json:"script" binding:"required"`
}

type NewServiceResponseStruct struct {
	Url   string `json:"url"`
	Error string `json:"error"`
}

// Custom validation function for HTTP/HTTPS scheme
func isHttpOrHttps(fl validator.FieldLevel) bool {
	uri := fl.Field().String()
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		return false
	}
	// Check if the scheme is either http or https
	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}

func (h ServiceController) ServiceProvisionner(c *gin.Context) {

	var payload NewServiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"url":   nil,
			"error": err.Error(), // will return EOF if body is empty
		})

	} else {

		// nameParam := c.Params.ByName("name")

		// if nameParam == "" {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"url":   nil,
		// 		"error": fmt.Errorf("Missing name parameter"),
		// 	})
		// }

		c.JSON(http.StatusOK, gin.H{
			"url":   "<DEPLOYED_URI>", // Deployed URL
			"error": nil,
		})

	}

}
