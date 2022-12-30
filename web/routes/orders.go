package routes

import (
	"net/http"

	"github.com/luminishion/wildberries-l0/orders"

	"github.com/gin-gonic/gin"
)

func getOrders(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		list := orders.List()

		c.HTML(http.StatusOK, "list.html", gin.H{
			"Orders": list,
		})

		return
	}

	order := orders.Get(id)
	if order == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "order not found",
		})

		return
	}

	c.Data(http.StatusOK, "application/json", order)
}

func Orders(r *gin.Engine) {
	r.GET("/orders", getOrders)
}
