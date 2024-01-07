package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yodeman/analyses-api/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"

	authorizationTypeToken = "bearer"
)

// authMiddleware ensures that requests carries authentication token.
// It also verifies the token carried by the request.
//
// Aborts a request if either authorization header is missing or token
// is invalid.
func authMiddleware(tokenMaker *token.PasetoMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("Invalid authorization header!")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		switch authorizationType {
		case authorizationTypeToken:
			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
				return
			}
			ctx.Set(authorizationPayloadKey, payload)
			ctx.Next()
		default:
			err := fmt.Errorf("Unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
		}
	}

}
