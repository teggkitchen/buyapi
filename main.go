package main

import (
	_ "buyapi/database"
	orm "buyapi/database"
	"buyapi/router"
)

func main() {
	defer orm.GormOpen.Close()
	router := router.InitRouter()
	router.Run(":8000")
}
