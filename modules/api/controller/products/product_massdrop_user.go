package products

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func (this API) MassdropUserAction(c *gin.Context) {

	slug := c.Param("id")

	if len(slug) < 1 {
		c.JSON(400, gin.H{"message": "Invalid request, need component slug.", "status": "error"})
		return
	}

	products := this.GCommerce.Products()
	product, err := products.GetByBson(bson.M{"slug": slug})

	if err != nil {
		c.JSON(404, gin.H{"message": "Invalid request, product not found.", "status": "error"})
		return
	}

	// Load Massdrop information (if exists)
	product.InitializeMassdrop()
	

	c.JSON(200, product)
}
