package webhandler

import (
	"context"
	"net/http"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	CreateAccountFeature interface {
		CreateAccount(context.Context, domain.User) error
	}
	createAccount struct {
		feature CreateAccountFeature
	}
)

func NewCreateAccount(f CreateAccountFeature) gin.HandlerFunc {
	handler := createAccount{f}
	handlerFunc := func(c *gin.Context) {
		user := User{}
		if err := c.ShouldBindBodyWithJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, BuildErrorResponse(err, BadRequestErrorCode))
			return
		}
		if err := ValidateAccountRequest(user); err != nil {
			logrus.WithContext(c).WithError(err).Error("Invalid request")
			c.AbortWithStatusJSON(http.StatusBadRequest, BuildErrorResponse(err, InvalidCredentialsErrorCode))
			return
		}

		err := handler.feature.CreateAccount(c, BuildUser(user))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, BuildErrorResponse(err, CreateUserErrorCode))
			return
		}

		c.Status(http.StatusCreated)
		return
	}
	return handlerFunc
}
