package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"oceanlearn.teach/ginessential/common"
	"oceanlearn.teach/ginessential/dto"
	"oceanlearn.teach/ginessential/model"
	"oceanlearn.teach/ginessential/response"
)

func Login(ctx *gin.Context){
	DB := common.GetDB()

	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 数据验证
	if len(telephone)!= 11{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password)<6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}

	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户不存在")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password)); err!=nil{
		response.Response(ctx,http.StatusBadRequest,400,nil,"密码错误")
	}

	// 发放token给前端
	token, err := common.ReleaseToken(user)
	if err!=nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,"系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	response.Success(ctx,gin.H{"token":token},"登陆成功")
}

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// 数据验证
	if len(telephone)!= 11{
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password)<6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}

	if len(name) == 0{
		name = "ss"
	}

	log.Println(name,telephone,password)


	// 判断手机号是否存在
	if isTelephoneExist(DB,telephone){
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户已经存在")
		return
	}

	// 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err !=nil {
		response.Response(ctx,http.StatusInternalServerError,500,nil,"加密错误")
		return
	}
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)

	// 返回结果
	response.Success(ctx,nil,"注册成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool{
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID !=0{
		return true
	}
	return false
}

func Info(ctx *gin.Context){
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code":200, "data": gin.H{"user":dto.ToUserDto(user.(model.User))}})
}