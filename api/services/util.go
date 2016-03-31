package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/middlewares"
	"github.com/mvader/pinlist/api/models"
)

func internalError(ctx *gin.Context, err error) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
	log.Err(err)
}

func ok(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func notFound(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
}

func unauthorized(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func badRequest(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusBadRequest)
}

func userFromCtx(c *gin.Context) *models.User {
	return c.MustGet(middlewares.UserKey).(*models.User)
}

func idFromCtx(c *gin.Context) int64 {
	n, _ := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	return n
}

const (
	defaultLimit = 20
	maxLimit     = 100

	limitParam  = "limit"
	offsetParam = "offset"
)

func limitAndOffset(ctx *gin.Context) (int, int64) {
	limit, _ := strconv.Atoi(ctx.Query(limitParam))
	if limit == 0 || limit > maxLimit {
		limit = defaultLimit
	}

	after, _ := strconv.ParseInt(ctx.Query(offsetParam), 10, 64)
	return limit, after
}

func intParam(ctx *gin.Context, param string) int {
	n, _ := strconv.Atoi(ctx.Query(param))
	return n
}

func intParamOrDefault(ctx *gin.Context, param string, def int) int {
	n := intParam(ctx, param)
	if n <= 0 {
		return def
	}
	return n
}
