package api

import (
	"context"
	"net/http"
	"oauth-server/app/model"
	"oauth-server/app/service"
	_errors "oauth-server/package/errors"
	"oauth-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type oAuthHandler struct {
	services service.ServiceCollections
}

func NewApiOAuthController(router *gin.Engine, services service.ServiceCollections) {
	handler := oAuthHandler{services}

	group := router.Group("api/v1/oauth")
	{
		group.POST("/refresh", handler.refreshToken)
	}
}

func (h *oAuthHandler) refreshToken(c *gin.Context) {
	var data model.RefreshTokenRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.OAuthSvc.RefreshToken(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
