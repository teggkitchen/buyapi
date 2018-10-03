# 會員購物車 For Golang Gin Restful API



## 運行專案前置

下載專案

將專案複製到 $GOPATH/src/buyapi


## 下載第三方套件 指令

下載第三方套件

```cmd
go get github.com/dgrijalva/jwt-go

go get github.com/gin-gonic/gin

go get github.com/go-sql-driver/mysql

go get github.com/jinzhu/gorm
```

## 創建資料庫

這裡使用Docker MySQL

docker指令：

```cmd
docker run --name buyapi_mysql -p 3307:3306 --net=aappii -v /Users/tingk/DockerProject/buyapi-mysql/mysql-data:/var/lib/mysql -v /Users/tingk/DockerProject/buyapi-mysql/mysql-config:/etc/mysql/conf.d -e MYSQL_ROOT_PASSWORD=12345600 -d mysql:8.0.12 --character-set-server=utf8 --collation-server=utf8_unicode_ci --init-connect='SET NAMES UTF8;'
```

## 資料庫結構

```cmd
Products Table
|---Id
|---名稱 
|---圖片 
|---價錢
|---創建時間
|---修改時間
```

```cmd
Orders Table
|---Id 
|---會員Id 
|---創建時間
|---修改時間
```

```cmd
OrderDetails Table
|---Id
|---訂單Id
|---商品Id
|---數量
```

```cmd
Members Table
|---Id
|---信箱
|---密碼
|---token
|---手機
|---是否驗證信箱
|---是否驗證手機
|---創建時間
|---修改時間
```

##創建資料庫步驟

創建資料庫：
```cmd
create database BUYDB character set utf8;
```

使用資料庫：
```cmd
use BUYDB;
```

創建Products資料表：
```cmd
CREATE TABLE Products(
id  INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
name VARCHAR(100),
img VARCHAR(100),
price INT,
created_at DATETIME DEFAULT NULL,
updated_at DATETIME DEFAULT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
```


創建Orders資料表：
```cmd
CREATE TABLE Orders(
id  INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
member_id INT,
created_at DATETIME DEFAULT NULL,
updated_at DATETIME DEFAULT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
```


創建OrderDetails資料表：
```cmd
CREATE TABLE OrderDetails(
id  INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
order_id INT,
product_id INT,
num INT
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

創建Members資料表：
```cmd
# CREATE TABLE Members(
# id  INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
# email VARCHAR(100),
# password VARCHAR(100),
# token VARCHAR(255),
# phone VARCHAR(100),
# is_email_verify INT(1),
# is_phone_verify INT(1),
# created_at DATETIME DEFAULT NULL,
# updated_at DATETIME DEFAULT NULL
# )ENGINE=InnoDB DEFAULT CHARSET=utf8;
```




## 運行Go

```cmd
go run main.go
```


## 備註

因一些環境參數的因素

1.建議使用go run main.go
2.docker不用network以免影響sql的連線
