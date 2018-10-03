package models

import (
	msg "buyapi/config"
	configDB "buyapi/database"
	"errors"
	"fmt"
	"time"
)

type RequestOrder struct {
	Token       string               `json:"token"`         // 會員憑證
	OrderDetail []RequestOrderDetail `json:"order_details"` // 訂單明細
}

type RequestOrderDetail struct {
	ProductId int64 `json:"product_id"` // 商品Id
	Num       int64 `json:"num"`        // 數量
}

type Order struct {
	Id        int64     `json:"id"`        // 訂單Id
	MemberId  int64     `json:"member_id"` // 會員Id
	CreatedAt time.Time `json:"createdAt"` // 開始時間
	UpdatedAt time.Time `json:"updatedAt"` // 更新時間
}

type OrderDetail struct {
	Id        int64 `json:"id"`         // 訂單明細Id
	OrderId   int64 `json:"order_id"`   // 訂單Id
	ProductId int64 `json:"product_id"` // 商品Id
	Num       int64 `json:"num"`        // 數量
}

// 查詢全部
// func (order *Order) QueryOrders(memberId int64) (data interface{}, err error) {
// 	var orders []Order
// 	result := configDB.GormOpen.Table("Orders").Where("member_id=?", memberId).Find(&orders)
// 	if result.Error != nil {
// 		err = result.Error
// 		return nil, err
// 	} else if len(orders[0:]) == 0 {
// 		return nil, errors.New(msg.NOT_FOUND_DATA_ERROR)
// 	}
// 	return orders, nil
// }

// 查詢訂單明細
func QueryOrderDetail(orderId int64) (data interface{}, err error) {
	var orderDetail []OrderDetail
	result := configDB.GormOpen.Table("OrderDetails").Where("order_id=?", orderId).Find(&orderDetail)
	if result.Error != nil {
		err = result.Error
		return nil, err
	} else if len(orderDetail[0:]) == 0 {
		return nil, errors.New(msg.NOT_FOUND_DATA_ERROR)
	}

	return orderDetail, nil
}

// 新增訂單
func (order *Order) InsertOrder(orderInfo Order, detailsInfo []OrderDetail) (data *Order, err error) {
	if err := configDB.GormOpen.Table("Orders").Create(&orderInfo).Error; err != nil {
		return nil, errors.New(msg.SQL_WRITE_ERROR)
	}

	// 設置訂單明細OrderId
	for i, _ := range detailsInfo {
		detailsInfo[i].OrderId = orderInfo.Id
	}

	//新增訂單明細
	if err := InsertOrderdetail(detailsInfo); err != nil {
		fmt.Println("332211")
		return nil, errors.New(msg.SQL_WRITE_ERROR)
	}

	return &orderInfo, nil
}

//新增訂單明細
func InsertOrderdetail(detailsInfo []OrderDetail) (err error) {
	sql := "INSERT INTO `OrderDetails` (`order_id`,`product_id`,`num`) VALUES "
	count := len(detailsInfo[0:]) - 1
	for i, detail := range detailsInfo {
		if i == (count) {
			sql += fmt.Sprintf("('%d','%d','%d');", detail.OrderId, detail.ProductId, detail.Num)
		} else {
			sql += fmt.Sprintf("('%d','%d','%d'),", detail.OrderId, detail.ProductId, detail.Num)
		}
	}
	fmt.Println(sql)

	tx := configDB.GormOpen.Begin()

	if err := configDB.GormOpen.Exec(sql).Error; err != nil {
		tx.Rollback()
		return errors.New(msg.SQL_WRITE_ERROR)
	}

	tx.Commit()
	return nil
}

// 刪除訂單
func (order *Order) Destroy(order_id int64) (err error) {
	var tmpOrder Order
	var tmpOrderDetail []OrderDetail
	if err = configDB.GormOpen.Table("Orders").Where("id=?", order_id).First(&tmpOrder, order_id).Error; err != nil {
		return err
	}

	if err = configDB.GormOpen.Table("OrderDetails").Where("order_id=?", order_id).Find(&tmpOrderDetail, order_id).Error; err != nil {
		return err
	}

	if err = configDB.GormOpen.Table("Orders").Where("id=?", order_id).Delete(&tmpOrder).Error; err != nil {
		return err
	}

	if err = configDB.GormOpen.Table("OrderDetails").Where("order_id=?", order_id).Delete(&tmpOrderDetail).Error; err != nil {
		return err
	}
	return nil
}

// 專用-查詢會員的全部訂單
type MemberOrder struct {
	OrderId   int64     `json:"order_id"`  // 訂單Id
	Num       int64     `json:"num"`       // 數量
	Price     int64     `json:"price"`     // 價錢
	Count     int64     `json:"count"`     // 價錢
	CreatedAt time.Time `json:"createdAt"` // 開始時間
}

//查詢會員的全部訂單
func (order *Order) QueryMemberOrder(token string) (memberOrder []MemberOrder, err error) {
	var tempCountOrder []MemberOrder

	sql := "select C.order_id, C.num, D.price, B.created_at "
	sql += "from OrderDetails as C "
	sql += "inner join Products as D on C.product_id = D.id "
	sql += "inner join Orders as B on B.id=C.order_id "
	sql += "inner join Members as A on A.id = B.member_id "
	sql += fmt.Sprintf("and A.token='%s';", token)
	fmt.Println(sql)
	if err := configDB.GormOpen.Raw(sql).Scan(&tempCountOrder).Error; err != nil {
		return nil, errors.New(msg.SQL_QUERY_ERROR)
	}

	for _, order := range tempCountOrder {
		if len(memberOrder) == 0 {
			// 第一筆
			memberOrder = append(memberOrder, order)
			memberOrder[0].Count = (order.Price * order.Num)
			fmt.Println("第一項queryCountOrder []--->", len(memberOrder)-1, memberOrder)

		} else {
			// 第一筆之後
			if order.OrderId == memberOrder[len(memberOrder)-1].OrderId {
				//重複 - 累加Count 不新增order
				fmt.Println("運行已重複")
				memberOrder[len(memberOrder)-1].Count += (order.Price * order.Num)
				fmt.Println("已重複queryCountOrder []--->", len(memberOrder)-1, memberOrder)
			} else {
				//不重複 - 新增order
				fmt.Println("運行未重複")
				memberOrder = append(memberOrder, order)
				memberOrder[len(memberOrder)-1].Count = (order.Price * order.Num)
				fmt.Println("未重複queryCountOrder []--->", len(memberOrder)-1, memberOrder)

				//清空暫存值
				memberOrder[len(memberOrder)-2].Num = 0
				memberOrder[len(memberOrder)-2].Price = 0
			}
		}

		fmt.Println("當前queryCountOrder", memberOrder)
		fmt.Println("當前queryCountOrder總數", len(memberOrder))
		fmt.Println("當前queryCountOrder總價錢", memberOrder[len(memberOrder)-1].Count)

	}

	if len(memberOrder)-1 > 0 {
		//最後一項 清空暫存值
		memberOrder[len(memberOrder)-1].Num = 0
		memberOrder[len(memberOrder)-1].Price = 0

	}
	fmt.Println("結果:", memberOrder)

	return memberOrder, nil
}

// 專用-查詢訂單明細
type MemberOrderDetail struct {
	OrderId   int64  `json:"order_id"`   // 訂單Id
	ProductId int64  `json:"product_id"` // 商品Id
	Name      string `json:"name"`       // 商品名稱
	Img       string `json:"img"`        // 圖片
	Num       int64  `json:"num"`        // 數量
	Price     int64  `json:"price"`      // 價錢
}

//查詢訂單明細
func (order *Order) QueryMemberOrderDetail(order_id int64) (memberOrderDetail []MemberOrderDetail, err error) {
	var tempCountOrderDetail []MemberOrderDetail

	sql := "select A.order_id, A.product_id, B.name, B.img, A.num, B.price "
	sql += "from OrderDetails as A "
	sql += "inner join Products as B on A.product_id = B.id "
	sql += fmt.Sprintf("and A.order_id='%d';", order_id)
	fmt.Println(sql)

	if err := configDB.GormOpen.Raw(sql).Scan(&tempCountOrderDetail).Error; err != nil {
		return nil, errors.New(msg.SQL_QUERY_ERROR)
	}

	memberOrderDetail = tempCountOrderDetail
	fmt.Println("結果:", memberOrderDetail)

	return memberOrderDetail, nil
}

// 檢查Product
func CheckProduct(orderDetails []RequestOrderDetail) (err error) {

	var tmpInt []int64

	fmt.Println("orderDetails:::", len(orderDetails))
	sql := "select id "
	sql += "from Products "
	for i := range orderDetails {
		if i == 0 {
			sql += fmt.Sprintf("where id in ('%d'", orderDetails[0].ProductId)

		} else if i > 0 {
			sql += fmt.Sprintf(",'%d'", orderDetails[i].ProductId)

		}
	}
	sql += ");"

	fmt.Println("sql: ", sql)

	if err := configDB.GormOpen.Raw(sql).Scan(&tmpInt).Error; err != nil {
		return errors.New(msg.SQL_QUERY_ERROR)
	}

	fmt.Println("tmpInt:::", len(tmpInt))

	if len(orderDetails) == len(tmpInt) {
		fmt.Println("無誤")
		return nil
	}
	return errors.New(msg.SQL_QUERY_ERROR)

}
