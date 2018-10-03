package apis

import (
	code "buyapi/config"
	msg "buyapi/config"
	model "buyapi/models"
	. "buyapi/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 會員註冊
func MemberSignUp(c *gin.Context) {
	var member model.Member

	// 取得參數
	member.Email = c.Request.FormValue("email")
	member.Phone = c.Request.FormValue("phone")
	member.Password = c.Request.FormValue("password")
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()
	member.Token = GetToken()

	// 參數是否有值
	if len(member.Email) > 0 && len(member.Phone) > 0 && len(member.Password) > 0 {
		// 驗證規則
		if IsEmail(member.Email) && IsPhone(member.Phone) {

			// 執行-增加會員
			result, err := member.Insert(member.Email)
			if err != nil {
				// 註冊失敗
				ShowJsonMSG(c, code.ERROR, msg.SIGNUP_ERROR)
				return
			}
			// 註冊成功
			ShowJsonDATA(c, code.SUCCESS, msg.SIGNUP_SUCCESS, result)
		} else {
			// 驗證失敗
			ShowJsonMSG(c, code.ERROR, msg.VERIFY_ERROR)
			return
		}

	} else {
		// 缺少參數
		ShowJsonMSG(c, code.ERROR, msg.ARGS_ERROR)
		return
	}

}

// 會員登入
func MemberSignIn(c *gin.Context) {
	var member model.Member

	// 取得參數
	member.Email = c.Request.FormValue("email")
	member.Password = c.Request.FormValue("password")

	// 參數是否有值
	if len(member.Email) > 0 && len(member.Password) > 0 {
		// 驗證規則
		if IsEmail(member.Email) {

			// 執行-增加會員
			result, err := member.Query(member.Email, member.Password)
			if err != nil {
				// 登入失敗
				ShowJsonMSG(c, code.ERROR, msg.SIGNIN_ERROR)
				return
			}

			ShowJsonDATA(c, code.SUCCESS, msg.SIGNIN_SUCCESS, result)
		} else {
			// 驗證失敗
			ShowJsonMSG(c, code.ERROR, msg.VERIFY_ERROR)
			return
		}

	} else {
		// 缺少參數
		ShowJsonMSG(c, code.ERROR, msg.ARGS_ERROR)
		return
	}

}
