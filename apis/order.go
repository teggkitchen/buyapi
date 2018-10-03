package apis

import (
	code "buyapi/config"
	msg "buyapi/config"
	model "buyapi/models"
	. "buyapi/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 展示全部訂單
func ShowOrders(c *gin.Context) {
	var order model.Order
	fmt.Println(len(strconv.Itoa(0)))
	fmt.Println(len(strconv.Itoa(1)))
	fmt.Println(len(strconv.Itoa(2)))
	fmt.Println(len(strconv.Itoa(01)))
	fmt.Println(len(strconv.Itoa(12)))
	token := c.Request.FormValue("token")

	//參數是否有值
	if len(token) > 0 {

		_, err := model.CheckToken(token)
		if err != nil {
			fmt.Println(msg.NOT_SIGNIN)
			ShowJsonMSG(c, code.ERROR, msg.NOT_SIGNIN)
			return
		}

		// 執行-查詢全部訂單
		result, err := order.QueryMemberOrder(token)
		fmt.Println(result)
		if err != nil {
			ShowJsonMSG(c, code.ERROR, msg.NOT_FOUND_DATA_ERROR)
			return
		}
		ShowJsonDATA(c, code.SUCCESS, msg.EXEC_SUCCESS, result)

	} else {
		// 缺少參數
		ShowJsonMSG(c, code.ERROR, msg.ARGS_ERROR)
		return
	}

}

// 展示一筆訂單明細
func ShowOrderDetail(c *gin.Context) {
	var order model.Order
	orderId := c.Request.FormValue("order_id")

	//參數是否有值
	if len(orderId) > 0 {
		order_id, _ := strconv.Atoi(orderId)

		// 執行-查詢一筆訂單明細
		result, err := order.QueryMemberOrderDetail(int64(order_id))

		if err != nil {
			ShowJsonMSG(c, code.ERROR, msg.NOT_FOUND_DATA_ERROR)
			return
		}
		ShowJsonDATA(c, code.SUCCESS, msg.EXEC_SUCCESS, result)

	} else {
		// 缺少參數
		ShowJsonMSG(c, code.ERROR, msg.ARGS_ERROR)
		return
	}

}

// 增加訂單
func CreateOrder(c *gin.Context) {

	var request model.RequestOrder

	// 取得回傳的JSON
	err := c.BindJSON(&request)
	if err != nil {
		//沒有資料
		fmt.Println(msg.REQUEST_DATA_ERROR)
	}

	var token string
	var order model.Order
	var orderDetails []model.OrderDetail

	token = request.Token
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	memberId, err := model.CheckToken(token)
	if err != nil {
		fmt.Println(msg.NOT_SIGNIN)
		ShowJsonMSG(c, code.ERROR, msg.NOT_SIGNIN)
		return
	}
	order.MemberId = memberId

	tmp := make([]model.OrderDetail, len(request.OrderDetail))

	for i, detail := range request.OrderDetail {
		tmp[i].Num = detail.Num
		tmp[i].ProductId = detail.ProductId
		fmt.Println(detail.ProductId)
	}
	orderDetails = tmp

	err = model.CheckProduct(request.OrderDetail)
	if err != nil {
		fmt.Println(msg.NOT_FOUND_PRODUCT_ERROR)
		ShowJsonMSG(c, code.ERROR, msg.NOT_FOUND_PRODUCT_ERROR)
		return
	}

	// 執行-新增訂單
	result, err := order.InsertOrder(order, orderDetails)
	fmt.Println(result) // id
	if err != nil {
		// 註冊失敗
		ShowJsonMSG(c, code.ERROR, msg.SIGNUP_ERROR)
		return
	}

	ShowJsonDATA(c, code.SUCCESS, msg.CREATE_SUCCESS, "")

}

func DeleteOrder(c *gin.Context) {
	var order model.Order
	order_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println("Delete order_id", order_id)

	// 執行-刪除訂單
	err = order.Destroy(order_id)
	if err != nil {
		//刪除失敗
		ShowJsonMSG(c, code.ERROR, msg.DELETE_ERROR)
		return
	}

	ShowJsonDATA(c, code.SUCCESS, msg.DELETE_SUCCESS, "")
}
