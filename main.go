package main
import (

	"net/http"

	"github.com/gin-gonic/gin"
)
func main(){
	router := gin.Default()
	router.GET("/health",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"code":0,
			"message":"ok",
			"data":gin.H{
				"status":"healthy",
			}})
		})
	if err := router.Run(":8080");err != nil{
	panic(err)
	}
	}
