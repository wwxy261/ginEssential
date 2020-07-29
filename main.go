package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"oceanlearn.teach/ginessential/common"
)


func main() {
	db := common.InitDB()
	defer db.Close()

	r := gin.Default()
	r = CollectionRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




