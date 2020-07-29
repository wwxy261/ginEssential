package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type User struct{
	gorm.Model
	Name string      `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password string  `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		// 数据验证
		if len(telephone)!= 11{
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须为11位"})
			return
		}
		if len(password)<6 {
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422, "msg":"密码不能小于6位"})
			return
		}

		if len(name) == 0{
			name = "ss"
		}

		log.Println(name,telephone,password)


		// 判断手机号是否存在
		if isTelephoneExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422, "msg":"用户已经存在"})
			return
		}

		// 创建用户
		newUser := User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser)


		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool{
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID !=0{
		return true
	}
	return false
}


func InitDB() *gorm.DB{
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "root1234"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName,args)
	if err != nil{
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}