package main

import (
	"github.com/gin-gonic/gin"
	"oceanlearn.teach/ginessential/controller"
)

func CollectionRoute(r *gin.Engine) *gin.Engine{
	r.POST("/api/auth/register", controller.Register)
	return r
}
