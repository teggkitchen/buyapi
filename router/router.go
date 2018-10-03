package router

import (
	. "buyapi/apis"
	"buyapi/config"
	// . "buyapi/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/product", CreateProduct)
	router.GET("/products", ShowProducts)
	router.PUT("/product/:id", UpdateProduct)
	router.DELETE("/product/:id", DestroyProduct)

	router.POST("/member/signup", MemberSignUp)
	router.POST("/member/signin", MemberSignIn)

	router.POST("/order/create", CreateOrder)
	router.POST("/order/query", ShowOrders)
	router.POST("/order/querydetail", ShowOrderDetail)
	router.DELETE("/order/delete/:id", DeleteOrder)

	//  ---訪問圖片---
	// go run 使用
	router.Static("/image", config.IMAGE_PATH)

	//  go build 使用
	// router.Static("/image", GetAppPath()+config.IMAGE_PATH2)

	return router
}
