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

type accountHandler struct {
	services service.ServiceCollections
}

func NewApiUserController(router *gin.Engine, services service.ServiceCollections) {
	handler := accountHandler{services}

	group := router.Group("api/v1/users")
	{
		group.POST("/login", handler.login)
		group.POST("/register", handler.register)
		group.POST("/logout", handler.logout)
	}
}

func (h *accountHandler) login(c *gin.Context) {
	var data model.LoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.UserSvc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *accountHandler) register(c *gin.Context) {
	var data model.RegisterRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.UserSvc.Register(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *accountHandler) logout(c *gin.Context) {
	var data model.LogoutRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	res, err := h.services.UserSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
