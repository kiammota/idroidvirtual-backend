package gateways

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](ctx *gin.Context, body *T) bool {
	if err := ctx.ShouldBindJSON(body); err != nil {
		Error(
			ctx,
			http.StatusBadRequest,
			"dados enviados em formato inválido",
		)
		return false
	}

	return true
}
