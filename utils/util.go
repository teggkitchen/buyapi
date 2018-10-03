package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//輸出失敗Json訊息
func ShowJsonMSG(c *gin.Context, code int64, msg string) {
	msg = strings.Replace(msg, "\b", "", -1)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}

//輸出成功Json訊息
func ShowJsonDATA(c *gin.Context, code int64, msg string, data interface{}) {
	msg = strings.Replace(msg, "\b", "", -1)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"data":    data,
		"message": msg,
	})
}

// 驗證Email
func IsEmail(email string) bool {
	reg := regexp.MustCompile(`^[_a-z0-9-]+([.][_a-z0-9-]+)*@[a-z0-9-]+([.][a-z0-9-]+)*$`)
	if m := reg.MatchString(email); m {
		return true
	} else {
		return false
	}
}

// 驗證手機號碼
func IsPhone(phone string) bool {
	reg := regexp.MustCompile(`^09[0-9]{8}$`)
	if m := reg.MatchString(phone); m {
		return true
	} else {
		return false
	}
}

//取得JWT Token
func GetToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"timeA": time.Now().Unix(),
		"timeB": time.Now(),
		"timeC": time.Now().Local(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("Token錯誤")

	}
	return tokenString
}

// MD5生成
func GetMD5Hash(text string) string {
	haser := md5.New()
	haser.Write([]byte(text))
	return hex.EncodeToString(haser.Sum(nil))
}

// 判斷資料夾是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat獲得資料夾訊息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 取得執行檔當前路徑，go run 無作用
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
