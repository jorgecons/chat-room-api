package webhandler

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
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
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := ValidateAccountRequest(user); err != nil {
			logrus.WithContext(c).WithError(err).Error("Invalid request")
			c.AbortWithStatusJSON(http.StatusBadRequest, BuildErrorResponse(InvalidUserCode, err))
			return
		}

		err := handler.feature.CreateAccount(c, BuildUser(user))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, BuildErrorResponse(InvalidUserCode, err))
			return
		}

		c.Status(http.StatusCreated)
		return
	}
	return handlerFunc
}
