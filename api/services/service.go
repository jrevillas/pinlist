package services

import "github.com/gin-gonic/gin"

// Service is the interface that defines an API service.
type Service interface {
	// Register is a method that, given a router, adds all the
	// routes needed by the service to that router.
	Register(*gin.RouterGroup)
}

// Services is a collection of services.
type Services []Service

// Register calls the Register method of all the services in
// the collection.
func (s Services) Register(r *gin.RouterGroup) {
	for _, service := range s {
		service.Register(r)
	}
}
