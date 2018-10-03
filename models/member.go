package models

import (
	msg "buyapi/config"
	configDB "buyapi/database"
	"errors"
	"fmt"
	"time"
)

type Member struct {
	Id            int64     `json:"id"`              // 會員Id
	Email         string    `json:"email"`           // 信箱
	Phone         string    `json:"phone"`           // 手機
	Password      string    `json:"password"`        // password
	Token         string    `json:"token"`           // token憑證
	IsEmailVerify int64     `json:"is_email_verify"` // 信箱是否驗證
	IsPhoneVerify int64     `json:"is_phone_verify"` // 手機是否驗證
	CreatedAt     time.Time `json:"createdAt"`       // 開始時間
	UpdatedAt     time.Time `json:"updatedAt"`       // 更新時間
}

type ShowToken struct {
	Token string `json:"token"` // token憑證
}

// 註冊會員
func (member *Member) Insert(email string) (data *Member, err error) {
	if isMemberReport(email) {
		return nil, errors.New(msg.SIGNUP_ERROR)
	}

	if err := configDB.GormOpen.Table("Members").Create(&member).Error; err != nil {
		return nil, errors.New(msg.SIGNUP_ERROR)
	}
	fmt.Println(member.Email)
	return member, nil
}

// 登入會員
func (member *Member) Query(email string, password string) (data *Member, err error) {
	if err := configDB.GormOpen.Table("Members").Where("email=? AND password=?", email, password).Scan(&member).Error; err != nil {
		return member, err
	}
	return member, nil
}

// 查詢重複會員,true：重複
func isMemberReport(email string) (b bool) {
	var memberVerify Member
	configDB.GormOpen.Table("Members").Where("email=?", email).Scan(&memberVerify)
	result := memberVerify.Email

	if len(result) > 0 {
		fmt.Println("2")
		return true
	}

	return false
}

// 檢查Token
func CheckToken(token string) (memberId int64, err error) {
	var member Member
	configDB.GormOpen.Table("Members").Where("token=?", token).Select([]string{"id , email"}).Scan(&member)
	memberId = member.Id
	if memberId <= 0 {
		return memberId, errors.New(msg.NOT_SIGNIN)
	}
	return memberId, nil
}
