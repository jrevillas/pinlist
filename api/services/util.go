package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/log"
)

func internalError(ctx *gin.Context, err error) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
	log.Err(err)
}
