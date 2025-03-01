package webhandler

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"chat-room-api/internal/core/domain"

	"github.com/gin-gonic/gin"
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
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := ValidateAccountRequest(user); err != nil {
			logrus.WithContext(c).WithError(err).Error("Invalid request")
			c.AbortWithStatusJSON(http.StatusBadRequest, BuildErrorResponse(InvalidUserCode, err))
			return
		}

		token, err := handler.feature.Login(c, BuildUser(user))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, BuildErrorResponse(InvalidUserCode, err))
			return
		}

		c.SecureJSON(http.StatusOK, Token{Token: token})
		return
	}
	return handlerFunc
}
