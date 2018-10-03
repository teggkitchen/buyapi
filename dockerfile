FROM golang:1.9.1

WORKDIR $GOPATH/src/buyapi

COPY ./ $GOPATH/src/buyapi

RUN go get github.com/dgrijalva/jwt-go

RUN go get github.com/gin-gonic/gin

RUN go get github.com/go-sql-driver/mysql

RUN go get github.com/jinzhu/gorm

RUN go build .

EXPOSE 8000

ENTRYPOINT  ["./buyapi"]
