package models

import (
	msg "buyapi/config"
	configDB "buyapi/database"
	"errors"
	"time"
)

type Product struct {
	Id        int64     `json:"id"`        // 商品Id
	Name      string    `json:"name"`      // 商品名稱
	Img       string    `json:"img"`       // 圖片
	Price     string    `json:"price"`     // 價錢
	CreatedAt time.Time `json:"createdAt"` // 開始時間
	UpdatedAt time.Time `json:"updatedAt"` // 更新時間

}

var Products []Product

// 新增商品
func (product *Product) Insert() (err error) {
	result := configDB.GormOpen.Table("Products").Create(&product)

	if result.Error != nil {
		err = result.Error
		return err
	}
	return nil
}

// 查詢全部
func (product *Product) QueryProducts() (data interface{}, err error) {
	var products []Product
	result := configDB.GormOpen.Table("Products").Find(&products)
	if result.Error != nil {
		err = result.Error
		return "", err
	} else if len(products[0:]) == 0 {
		return nil, errors.New(msg.NOT_FOUND_DATA_ERROR)
	}
	return products, nil
}

// 修改商品
func (product *Product) Update(id int64) (err error) {
	var updateProduct Product
	if err = configDB.GormOpen.Table("Products").Select([]string{"id"}).First(&updateProduct, id).Error; err != nil {
		return err
	}
	if err = configDB.GormOpen.Table("Products").Model(&updateProduct).Updates(&product).Error; err != nil {
		return err
	}
	return nil
}

// 刪除商品
func (product *Product) Destroy(id int64) (err error) {

	if err = configDB.GormOpen.Table("Products").Select([]string{"id"}).First(&product, id).Error; err != nil {
		return err
	}

	if err = configDB.GormOpen.Table("Products").Delete(&product).Error; err != nil {
		return err
	}
	return nil
}

// 查詢圖片名稱 - 指定id
func (product *Product) GetProductImg(id int64) (imgName string, err error) {

	var productRow Product
	configDB.GormOpen.Table("Products").Where("id=?", id).Scan(&productRow)

	result := productRow.Img
	if len(result) > 0 {
		return result, nil
	}

	return result, errors.New(msg.CONTINUE_NOT_FOUND_IMAGE)

}
