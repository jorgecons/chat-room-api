package webhandler

import (
	"context"
	"errors"
	"net/http"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type (
	LoginFeature interface {
		Login(context.Context, domain.User) (string, error)
	}
	login struct {
		feature LoginFeature
	}
)

func NewLogin(f LoginFeature) gin.HandlerFunc {
	handler := login{f}
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

		token, err := handler.feature.Login(c, BuildUser(user))
		if err != nil {
			if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrInvalidUser) {
				c.AbortWithStatusJSON(http.StatusBadRequest, BuildErrorResponse(ErrLoginError, InvalidCredentialsErrorCode))
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, BuildErrorResponse(err, InvalidCredentialsErrorCode))
			return
		}

		c.SecureJSON(http.StatusOK, Token{Token: token})
		return
	}
	return handlerFunc
}
