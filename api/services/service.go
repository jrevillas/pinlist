package services

import "github.com/gin-gonic/gin"

type Service interface {
	Register(*gin.RouterGroup)
}

type Services []Service

func (s Services) Register(r *gin.RouterGroup) {
	for _, service := range s {
		service.Register(r)
	}
}
